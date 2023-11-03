package ios

import "fmt"

func parsePort(v []string) (PortMatcher, []string, error) {
	switch v[0] {
	case "eq":
		t, remaining, err := parseTransportProtocol(v[1:])
		if err != nil {
			return PortMatch{}, v, err
		}
		return PortMatch{name: "eq", proto: t}, remaining, nil
	case "gt":
		t, remaining, err := parseTransportProtocol(v[1:])
		if err != nil {
			return PortMatch{}, v, err
		}
		return PortMatch{name: "gt", proto: t}, remaining, nil
	case "lt":
		t, remaining, err := parseTransportProtocol(v[1:])
		if err != nil {
			return PortMatch{}, v, err
		}
		return PortMatch{name: "lt", proto: t}, remaining, nil
	case "range":
		t1, remaining, err := parseTransportProtocol(v[1:])
		if err != nil {
			return PortMatch{}, v, err
		}
		t2, remaining, err := parseTransportProtocol(remaining)
		if err != nil {
			return PortMatch{}, v, err
		}
		return RangeMatch{name: "range", proto1: t1, proto2: t2}, remaining, err
	default:
		// no port match in the given string list
		// so return nothing and the string list unchanged
		return PortMatch{}, v, nil
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

type RangeMatch struct {
	name   string
	proto1 TransportProtocol
	proto2 TransportProtocol
}

func (r RangeMatch) Match(o PortMatcher) bool {
	return false
}

func (r RangeMatch) String() string {
	return fmt.Sprintf("range %s %s", r.proto1, r.proto2)
}
