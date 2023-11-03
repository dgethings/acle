package ios

import (
	"fmt"
	"strconv"
)

func parseTransportProtocol(v []string) (TransportProtocol, []string, error) {
	switch v[0] {
	case "bgp":
		return TransProto{"bgp", 179}, v[1:], nil
	case "179":
		return TransProto{"bgp", 179}, v[1:], nil
	case "ftp":
		return TransProto{"ftp", 21}, v[1:], nil
	case "ftp-data":
		return TransProto{"ftp-data", 20}, v[1:], nil
	default:
		p, err := strconv.Atoi(v[0])
		if err != nil {
			return TransProto{}, v, fmt.Errorf("unrecognised protocol %s", v[0])
		}
		return TransProto{v[0], uint16(p)}, v[1:], nil
	}
}

type TransportProtocol interface {
	String() string
	Integer() uint16
	Compare(TransProto) bool
}

type TransProto struct {
	name   string
	number uint16
}

func (t TransProto) String() string {
	return t.name
}

func (t TransProto) Integer() uint16 {
	return t.number
}

func (t TransProto) Compare(o TransProto) bool {
	return t.number == o.number
}
