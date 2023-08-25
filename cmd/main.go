package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

func readFile(f string) string {
	b, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("Failed to read %s, got %v\n", f, err)
	}
	return string(b)
}

func readStdin() string {
	stdin, err := io.ReadAll(os.Stdin)

	if err != nil {
		log.Fatalf("Failed to read stdin, got %v", err)
	}
	return string(stdin)
}

func main() {
	var cfg string
	var f string
	flag.StringVar(&f, "if", "", "Path to input config file")
	var acl_id string
	flag.StringVar(&acl_id, "acl_id", "", "Name or number of ACL")
	var acl_type string
	flag.StringVar(&acl_type, "acl_type", "standard", "Type of ACL. Either 'standard' or 'extended', default is 'standard'")
	flag.Parse()
	if f != "" {
		cfg = readFile(f)
	} else {
		cfg = readStdin()
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
