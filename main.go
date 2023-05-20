package main

import (
	"fmt"
)

func main() {
	fmt.Println(encodeDNSFlagsAsUint16(DNSFlags{isResponse: true}))
}
