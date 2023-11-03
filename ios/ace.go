package ios

import "fmt"

type ACE struct {
	Index      uint8
	Action     Action
	Protocol   IPProtocol
	SrcPrefix  IPNetwork
	SrcMatch   PortMatcher
	DestPrefix IPNetwork
	DestMatch  PortMatcher
}

func (ace ACE) String() string {
	return fmt.Sprintf("%s %s %s", ace.Action, ace.Protocol.String(), ace.SrcPrefix.String())
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
