package ios

import (
	"errors"
	"fmt"
)

type Action int8

const (
	NonAction Action = iota
	Permit
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

func parseAction(s string) (Action, error) {
	var action Action
	var err error
	switch s {
	case "permit":
		action = Permit
	case "deny":
		action = Deny
	default:
		action = NonAction
		msg := fmt.Sprintf("Unrecognised action %s", s)
		err = errors.New(msg)
	}
	return action, err
}
