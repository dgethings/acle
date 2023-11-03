package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	ios "github.com/dgethings/acle/ios"
)

func readFile(f string) string {
	b, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("Failed to read %s, got %v\n", f, err)
	}
	return string(b[:])
}

func readStdin() string {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read stdin, got %v", err)
	}
	return string(stdin)
}

func main() {
	var f string
	flag.StringVar(&f, "if", "", "Path to input config file")
	var acl_id string
	flag.StringVar(&acl_id, "acl_id", "", "Name or number of ACL")
	var acl_type string
	flag.StringVar(&acl_type, "acl_type", "standard", "Type of ACL. Either 'standard' or 'extended', default is 'standard'")
	flag.Parse()
	var cfg string
	if f != "" {
		cfg = readFile(f)
	} else {
		cfg = readStdin()
	}
	s, err := ios.GetACL(acl_id, cfg, acl_type)
	if err != nil {
		log.Fatal(err)
		return
	}
	acl, err := ios.NewACL(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(acl)
}
