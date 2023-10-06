package ios

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

func GetACL(name string, cfg string, acl_type string) ([]string, error) {
	var acl []string
	substr := fmt.Sprintf("access-list %s", name)
	for _, l := range strings.Split(cfg, "\n") {
		if strings.Contains(l, substr) {
			acl = append(acl, l)
		}
	}
	if len(acl) == 0 {
		return acl, fmt.Errorf("No such ACL named %s found", name)
	}
	return acl, nil
}

type ACL struct {
	id  int8
	ace []ACE
}

func NewACL(s []string) (ACL, error) {
	acl := ACL{}
	for i, ace := range s {
		msg := fmt.Errorf("Failed to parse ACE '%s'", ace)
		v := strings.Split(ace, " ")
		id, err := strconv.ParseInt(v[1], 10, 8)
		if err != nil {
			return acl, fmt.Errorf("Failed to parse ACL ID from %s, got %v", ace, err)
		}
		acl.id = int8(id)
		action, err := parseAction(v[2])
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		proto := parseIPProtocol(v[3])
		srcAddr, remaining, err := parseAddr(v[4:])
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		srcMatch, srcPort, remaining, err := parsePort(remaining)
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		dstAddr, remaining, err := parseAddr(remaining)
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		a := ACE{int8(i), action, proto, srcAddr, srcMatch, srcPort, dstAddr}
		acl.ace = append(acl.ace, a)
	}
	return acl, nil
}

func parseIPProtocol(s string) IPProtocol {
	var proto IPProtocol
	switch s {
	case IP.name:
		proto = IP
	case ICMP.name:
		proto = ICMP
	case UDP.name:
		proto = UDP
	case ESP.name:
		proto = ESP
	}
	return proto
}

func (acl ACL) String() string {
	var s []string
	for _, ace := range acl.ace {
		s = append(s, fmt.Sprintf("access-list %d %s", acl.id, ace))
	}
	return strings.Join(s, "\n")
}

type ACE struct {
	Index      int8
	Action     Action
	Protocol   IPProtocol
	SrcPrefix  IPNetwork
	SrcMatch   PortMatcher
	SrcPort    TransportProtocol
	DestPrefix IPNetwork
	// DestPort   TransportProtocol
}

func (ace ACE) String() string {
	return fmt.Sprintf("%s %s %s", ace.Action, ace.Protocol.String(), ace.SrcPrefix.String())
}

type Stringer interface {
	String() string
}
type Integerer interface {
	Integer() int8
}
type Equator interface {
	Equals() bool
}
type Gter interface {
	Gt() bool
}
type Lter interface {
	Lt() bool
}

type IPProtocol struct {
	name   string
	number int8
}

func (p IPProtocol) String() string {
	return p.name
}

func (p IPProtocol) Integer() int8 {
	return p.number
}

func (p IPProtocol) Equals(o IPProtocol) bool {
	return p.number == o.number
}

func (p IPProtocol) Gter(o IPProtocol) bool {
	return p.number > o.number
}

func (p IPProtocol) Lter(o IPProtocol) bool {
	return p.number < o.number
}

var IP = IPProtocol{"ip", 0}
var ICMP = IPProtocol{"icmp", 1}
var UDP = IPProtocol{"udp", 17}
var ESP = IPProtocol{"esp", 50}

type PortMatch struct {
	name string
}

type PortMatcher interface {
	Match(TransportProtocol) bool
	String() string
}

func (p PortMatch) String() string {
	return p.name
}

type PortEqual struct {
	name string
}

func (p PortEqual) Match(o TransportProtocol) bool {
	return false
}

func (p PortEqual) String() string {
	return p.name
}

func parsePort(v []string) (PortMatcher, TransportProtocol, []string, error) {
	var p PortMatcher
	var t TransportProtocol
	var remaining []string
	var err error
	switch v[0] {
	case "eq":
		p = PortEqual{"eq"}
		t, remaining, err = parseTransportProtocol(v[1:])
	}
	if err != nil {
		return p, t, remaining, err
	}
	return p, t, remaining, nil
}

func parseTransportProtocol(v []string) (TransportProtocol, []string, error) {
	switch v[0] {
	case "bgp":
		return BGP, v[1:], nil
	default:
		return TransProto{}, v, fmt.Errorf("Unrecognised protocol %s", v[0])
	}
}

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
	var src IPNetwork
	var msg error
	var remaining []string
	switch v[0] {
	case "any":
		net, err := netip.ParsePrefix("0.0.0.0/0")
		if err != nil {
			msg = fmt.Errorf("Somehow 0.0.0.0 is not a valid IP: %w", err)
		}
		src = IPNetwork{ip: net, isHost: false, isAny: true}
		remaining = v[1:]
	case "host":
		ip, err := netip.ParseAddr(v[1])
		if err != nil {
			msg = fmt.Errorf("'%s' is invalid host IP: %w", v[1], err)
		}
		net := netip.PrefixFrom(ip, 32)
		src = IPNetwork{ip: net, isHost: true, isAny: false}
		remaining = v[2:]
	default:
		ip, err := netip.ParseAddr(v[0])
		if err != nil {
			msg = fmt.Errorf("'%s' is not a valid IP address: %w", v[0], err)
		}
		prefixLen, err := prefixFromWildcard(v[1])
		if err != nil {
			msg = fmt.Errorf("%w", err)
		}
		net := netip.PrefixFrom(ip, prefixLen)
		src = IPNetwork{ip: net, isHost: false, isAny: false}
		remaining = v[2:]
	}
	return src, remaining, msg
}

func prefixFromWildcard(s string) (int, error) {
	octets := strings.Split(s, ".")
	// number of bits set to 1 in wildcard
	bits := 0
	var msg error
	fail := false
	// if len(octets) != 4 {
	// 	msg = fmt.Sprintf("'%s' could not be converted into 4 octets", s)
	// 	return -1, errors.New(msg)
	// }
	for i, octet := range octets {
		num, err := strconv.Atoi(octet)
		if err != nil {
			fail = true
			msg = fmt.Errorf("unable to parse octect %d of wildcard: %w", i, err)
		}
		if num > 0 {
			bits += (num + 1) / 32
		}
	}
	if fail {
		return -1, msg
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

type TransportProtocol interface {
	String() string
	Integer() uint8
	Compare(TransProto) bool
}

type TransProto struct {
	name   string
	number uint8
}

func (t TransProto) String() string {
	return t.name
}

func (t TransProto) Integer() uint8 {
	return t.number
}

func (t TransProto) Compare(o TransProto) bool {
	return t.number == o.number
}

var BGP = TransProto{"bgp", 179}
