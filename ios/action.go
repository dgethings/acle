package ios

import (
	"fmt"
)

type Action bool

func (a Action) String() string {
	if a {
		return "permit"
	} else {
		return "deny"
	}
}

func parseAction(s string) (Action, error) {
	var action Action
	switch s {
	case "permit":
		return true, nil
	case "deny":
		return false, nil
	case "":
		return action, fmt.Errorf("no action field found")
	default:
		return action, fmt.Errorf("unrecognised action %s", s)
	}
}
