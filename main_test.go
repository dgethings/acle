package main

import (
	"fmt"
	"testing"
)

func TestEquals(t *testing.T) {
	proto1 := Protocol{name: "Foo", number: 1}
	proto2 := Protocol{name: "Bar", number: 2}
	if !proto1.equals(proto1) {
		t.Errorf("%v fails to match itself", proto1)
	}
	if proto1.equals(proto2) {
		t.Errorf("%v should not match %v", proto1, proto2)
	}
}

func TestIPEqualsAllProtocols(t *testing.T) {
	IP := Protocol{name: "IP", number: 0}
	UDP := Protocol{name: "UDP", number: 7}
	if !IP.equals(UDP) {
		t.Errorf("%v should match %v", IP, UDP)
	}
}

func TestProtocolAsString(t *testing.T) {
	proto := Protocol{name: "Foo", number: 1}
	if fmt.Sprintf("%s", proto) == "Foo" {
		t.Errorf("printing %s does not render 'Foo'", proto)
	}
}

// func TestProtocolAsInt(t *testing.T) {
// 	proto := Protocol{name: "Foo", number: 1}
// 	if proto == 1 {
// 		t.Errorf("Protocol as integer %d does not match 1", proto)
// 	}
// }
