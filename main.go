package main

import "fmt"

const DOMAIN_NAME_SIMPLE_CASE = "google.com"
const DOMAIN_NAME_NEED_TO_RESOLVE_NAME_SERVER_IP = "twitter.com"
const DOMAIN_NAME_WITH_CNAME = "www.facebook.com"

func main() {
	fmt.Println("Welcome to our little DNS resolving showcase !")
	fmt.Println("\nFirst case is the most simple, we just iterate through the multiple DNS server")
	resolveDomainNameToIp(DOMAIN_NAME_SIMPLE_CASE)
	fmt.Print("\nSecond case is a bit more complicated : a DNS server IP is missing, so we need to resolve it too")
	resolveDomainNameToIp(DOMAIN_NAME_NEED_TO_RESOLVE_NAME_SERVER_IP)
	fmt.Println("\nThird case involves a CNAME record !")
	resolveDomainNameToIp(DOMAIN_NAME_WITH_CNAME)
}
