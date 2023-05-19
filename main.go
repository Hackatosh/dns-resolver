package main

import (
	"fmt"
)

func main() {
	fmt.Println(buildDNSFlagsAsBytes(DNSFlags{isResponse: true}))
}
