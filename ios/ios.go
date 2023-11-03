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
		return acl, fmt.Errorf("no such ACL named %s found", name)
	}
	return acl, nil
}

type ACL struct {
	ace []ACE
	id  int8
}

func NewACL(s []string) (ACL, error) {
	acl := ACL{}
	for i, ace := range s {
		msg := fmt.Errorf("failed to parse ACE '%s'", ace)
		v := strings.Split(ace, " ")
		id, err := strconv.ParseInt(v[1], 10, 8)
		if err != nil {
			return acl, fmt.Errorf("failed to parse ACL ID from %s, got %v", ace, err)
		}
		acl.id = int8(id)
		action, err := parseAction(v[2])
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		proto, err := parseIPProtocol(v[3])
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		srcAddr, remaining, err := parseAddr(v[4:])
		if err != nil {
			return acl, errors.Join(msg, fmt.Errorf("unable to parse %s", v[4:]), err)
		}
		srcMatch, remaining, err := parsePort(remaining)
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		dstAddr, remaining, err := parseAddr(remaining)
		if err != nil {
			return acl, errors.Join(msg, fmt.Errorf("unable to parse %s", remaining), err)
		}
		dstMatch, remaining, err := parsePort(remaining)
		if err != nil {
			return acl, errors.Join(msg, err)
		}
		if len(remaining) > 0 {
			return acl, errors.Join(msg, fmt.Errorf("following parts of the ACE have not been parsed: %s", remaining))
		}
		a := ACE{uint8(i), action, proto, srcAddr, srcMatch, dstAddr, dstMatch}
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
