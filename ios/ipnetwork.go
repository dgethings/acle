package ios

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

type IPNetwork struct {
	ip     netip.Prefix
	isHost bool
	isAny  bool
}

func (ip IPNetwork) String() string {
	if ip.isAny {
		return "any"
	}
	if ip.isHost {
		return fmt.Sprintf("host %s", ip.ip.Addr())
	}
	return fmt.Sprintf("%s %s", ip.ip.Addr(), wildcardFromPrefix(ip.ip.Bits()))
}

func parseAddr(v []string) (IPNetwork, []string, error) {
	var net IPNetwork
	var msg error
	if len(v) == 0 {
		return net, v, fmt.Errorf("given empty string: %s", v)
	}
	switch v[0] {
	case "any":
		prefix, err := netip.ParsePrefix("0.0.0.0/0")
		if err != nil {
			return net, v, fmt.Errorf("somehow 0.0.0.0 is not a valid IP: %w", err)
		}
		return IPNetwork{ip: prefix, isHost: false, isAny: true}, v[1:], msg
	case "host":
		ip, err := netip.ParseAddr(v[1])
		if err != nil {
			return net, v, fmt.Errorf("'%s' is invalid host IP: %w", v[1], err)
		}
		net := netip.PrefixFrom(ip, 32)
		return IPNetwork{ip: net, isHost: true, isAny: false}, v[2:], msg
	default:
		ip, err := netip.ParseAddr(v[0])
		if err != nil {
			return net, v, fmt.Errorf("'%s' is not a valid IP address: %w", v[0], err)
		}
		prefixLen, err := prefixFromWildcard(v[1])
		if err != nil {
			return net, v, fmt.Errorf("%w", err)
		}
		net := netip.PrefixFrom(ip, prefixLen)
		return IPNetwork{ip: net, isHost: false, isAny: false}, v[2:], msg
	}
}

func prefixFromWildcard(s string) (int, error) {
	octets := strings.Split(s, ".")
	// number of bits set to 1 in wildcard
	bits := 0
	// if len(octets) != 4 {
	// 	msg = fmt.Sprintf("'%s' could not be converted into 4 octets", s)
	// 	return -1, errors.New(msg)
	// }
	for i, octet := range octets {
		num, err := strconv.Atoi(octet)
		if err != nil {
			return -1, fmt.Errorf("unable to parse octect %d of wildcard: %w", i, err)
		}
		if num > 0 {
			bits += (num + 1) / 32
		}
	}
	return 32 - bits, nil
}

func wildcardFromPrefix(i int) string {
	oct := make([]int, 4)
	l := 32 - i
	for j := 3; j > -1; j-- {
		var o int
		if l >= 8 {
			o = 8
			l = l - 8
		} else {
			o = l
			l = 0
		}
		if o > 0 {
			oct[j] = (o * 32) - 1
		} else {
			oct[j] = 0
		}
	}
	return fmt.Sprintf("%d.%d.%d.%d", oct[0], oct[1], oct[2], oct[3])
}
