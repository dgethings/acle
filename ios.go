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
		var action Action
		switch v[2] {
		case "permit":
			action = Permit
		case "deny":
			action = Deny
		}
		a := ACE{int8(i), action}
		acl.ace = append(acl.ace, a)
	}
	return acl, nil
}

func (acl ACL) String() string {
	var s []string
	for _, ace := range acl.ace {
		s = append(s, fmt.Sprintf("access-list %d %s", acl.id, ace))
	}
	return strings.Join(s, "\n")
}

type ACE struct {
	Index  int8
	Action Action
	// Protocol   IPProtocol
	// SrcPrefix  IPNetwork
	// SrcPort    TransportProtocol
	// DestPrefix IPNetwork
	// DestPort   TransportProtocol
}

func (ace ACE) String() string {
	return fmt.Sprintf("%s", ace.Action)
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

type IPProtocol struct{}

type IPNetwork struct{}

type TransportProtocol struct{}
