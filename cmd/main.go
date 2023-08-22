package main

import (
	"flag"
	"fmt"

	ios "github.com/dgethings/acle"
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
	f := flag.String("cfg", "", "Path to config file")
	acl_id := flag.String("acl_id", "", "Name or number of ACL")
	acl_type := flag.String("acl_type", "standard", "Type of ACL. Either 'standard' or 'extended', default is 'standard'")
	flag.Parse()
	cfg, err := ios.LoadConfig(*f)
	if err != nil {
		fmt.Printf("%v\n", err)
		flag.Usage()
		return
	}
	acl, err := ios.GetACL(acl_id, cfg, acl_type)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, l := range acl {
		fmt.Println(l)
	}
}
