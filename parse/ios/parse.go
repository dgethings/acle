package ios

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func LoadConfig(f string) ([]byte, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		msg := fmt.Sprintf("Failed to read %s, got %v\n", f, err)
		return nil, errors.New(msg)
	}
	return b, nil
}

func GetACL(name string, cfg []byte) ([]string, error) {
	c := string(cfg[:])
	var acl []string
	substr := fmt.Sprintf("access-list %s", name)
	for _, l := range strings.Split(c, "\n") {
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
	ace []ACE
}

type ACE struct {
	Index      int
	Action     Action
	Protocol   IPProtocol
	SrcPrefix  IPNetwork
	SrcPort    TransportProtocol
	DestPrefix IPNetwork
	DestPort   TransportProtocol
}

type Action int8

const (
	Permit Action = iota
	Deny
)

func (a Action) String() string {
	switch a {
	case Permit:
		return "permit"
	case Deny:
		return "deny"
	}
}
