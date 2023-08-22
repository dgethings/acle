package main

import (
	"fmt"

	"github.com/dgethings/acle/parse/ios"
)

type Protocol struct {
	number int
	name   string
}

func (p Protocol) String() string {
	return p.name
}

func (p Protocol) Int() int {
	return p.number
}

func (p Protocol) equals(other Protocol) bool {
	if (p.number == 0) || (other.number == 0) {
		return true
	}
	return p.number == other.number
}

func main() {
	f := "./test_data/sample.ios"
	acl_id := "104"
	cfg, err := ios.LoadConfig(f)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	acl, err := ios.GetACL(acl_id, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, l := range acl {
		fmt.Println(l)
	}
}
