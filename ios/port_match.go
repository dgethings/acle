package ios

import "fmt"

func parsePort(v []string) (PortMatcher, []string, error) {
	var t TransportProtocol
	var remaining []string
	var err error
	switch v[0] {
	case "eq":
		t, remaining, err = parseTransportProtocol(v[1:])
		if err != nil {
			return PortMatch{}, v, err
		}
		return PortMatch{name: "eq", proto: t}, remaining, nil
	default:
		return PortMatch{}, v, err
	}
}

type PortMatcher interface {
	Match(PortMatcher) bool
	String() string
}

type PortMatch struct {
	name  string
	proto TransportProtocol
}

func (p PortMatch) Match(o PortMatcher) bool {
	return false
}

func (p PortMatch) String() string {
	if len(p.name) == 0 {
		return ""
	}
	return fmt.Sprintf("%s %s", p.name, p.proto)
}
