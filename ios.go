package ios

import (
	"errors"
	"fmt"
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
		v := strings.Split(ace, " ")
		id, err := strconv.ParseInt(v[1], 10, 8)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse ACL ID from %s, got %v", ace, err)
			return acl, errors.New(msg)
		}
		acl.id = int8(id)
		action := parseAction(v[2])
		proto := parseIPProtocol(v[3])
		a := ACE{int8(i), action, proto}
		acl.ace = append(acl.ace, a)
	}
	return acl, nil
}

func parseAction(s string) Action {
	var action Action
	switch s {
	case "permit":
		action = Permit
	case "deny":
		action = Deny
	}
	return action
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
	Index    int8
	Action   Action
	Protocol IPProtocol
	// SrcPrefix  IPNetwork
	// SrcPort    TransportProtocol
	// DestPrefix IPNetwork
	// DestPort   TransportProtocol
}

func (ace ACE) String() string {
	return fmt.Sprintf("%s %s", ace.Action, ace.Protocol.String())
}

type Action int8

const (
	Permit Action = iota
	Deny
)

func (a Action) String() string {
	var action string
	switch a {
	case Permit:
		action = "permit"
	case Deny:
		action = "deny"
	}
	return action
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

func (p *IPProtocol) String() string {
	return p.name
}

func (p *IPProtocol) Integer() int8 {
	return p.number
}

func (p *IPProtocol) Equals(o IPProtocol) bool {
	return p.number == o.number
}

func (p *IPProtocol) Gter(o IPProtocol) bool {
	return p.number > o.number
}

func (p *IPProtocol) Lter(o IPProtocol) bool {
	return p.number < o.number
}

var IP = IPProtocol{"ip", 0}
var ICMP = IPProtocol{"icmp", 1}
var UDP = IPProtocol{"udp", 17}
var ESP = IPProtocol{"esp", 50}

type IPNetwork struct{}

type TransportProtocol struct{}
