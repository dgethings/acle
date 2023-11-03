package ios

import "fmt"

func parseIPProtocol(s string) (IPProtocol, error) {
	switch s {
	case IP.name:
		return IP, nil
	case ICMP.name:
		return ICMP, nil
	case UDP.name:
		return UDP, nil
	case ESP.name:
		return ESP, nil
	default:
		return IPProtocol{}, fmt.Errorf("%s is not a known IP protocol", s)
	}
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

var (
	IP   = IPProtocol{"ip", 0}
	ICMP = IPProtocol{"icmp", 1}
	UDP  = IPProtocol{"udp", 17}
	ESP  = IPProtocol{"esp", 50}
)
