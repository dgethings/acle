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
		msg := fmt.Sprintf("No such ACL named %s found", name)
		return nil, errors.New(msg)
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
		msg := fmt.Sprintf("Failed to parse ACE '%s'", ace)
		v := strings.Split(ace, " ")
		id, err := strconv.ParseInt(v[1], 10, 8)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse ACL ID from %s, got %v", ace, err)
			return acl, errors.New(msg)
		}
		acl.id = int8(id)
		action, err := parseAction(v[2])
		if err != nil {
			return acl, errors.Join(errors.New(msg), err)
		}
		proto := parseIPProtocol(v[3])
		srcAddr, remaining, err := parseAddr(v[4:])
		if err != nil {
			return acl, errors.Join(errors.New(msg), err)
		}
		srcMatch, srcPort, remaining, err := parsePort(remaining)
		if err != nil {
			return acl, errors.Join(errors.New(msg), err)
		}
		dstAddr, remaining, err := parseAddr(remaining)
		if err != nil {
			return acl, errors.Join(errors.New(msg), err)
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
	return TransportProtocol{"foo"}, v, nil
}

type IPNetwork struct {
	ip     netip.Addr
	isHost bool
	isAny  bool
}

func (ip IPNetwork) String() string {
	if ip.isAny {
		return "any"
	}
	if ip.isHost {
		return fmt.Sprintf("host %s", ip.ip)
	}
	return fmt.Sprintf("%s 0.0.0.255", ip.ip)
}

func parseAddr(v []string) (IPNetwork, []string, error) {
	var src IPNetwork
	var msg string
	var remaining []string
	switch v[0] {
	case "any":
		ip, err := netip.ParseAddr("0.0.0.0")
		if err != nil {
			msg = fmt.Sprintf("Somehow 0.0.0.0 is not a valid IP: %v", err)
		}
		src = IPNetwork{ip: ip, isHost: false, isAny: true}
		remaining = v[1:]
	case "host":
		ip, err := netip.ParseAddr(v[1])
		if err != nil {
			msg = fmt.Sprintf("%s invalid host IP: %v", v[1], err)
		}
		src = IPNetwork{ip: ip, isHost: true, isAny: false}
		remaining = v[2:]
	default:
		ip, err := netip.ParseAddr(v[0])
		if err != nil {
			msg = fmt.Sprintf("%s is not a valid IP address: %v", v[0], err)
		}
		src = IPNetwork{ip: ip, isHost: false, isAny: false}
		remaining = v[2:]
	}
	if src.ip.IsValid() {
		return src, remaining, nil
	}
	return src, remaining, errors.New(msg)
}

type TransportProtocol struct {
	name string
}
