package main

import (
	"errors"
	"fmt"
)

// https://www.iana.org/domains/root/servers
const ROOT_SERVER_IP_ADDRESS = "198.41.0.4"

func sendQueryForDomainName(dnsServerIpAddress string, domainName string) DNSPacket {
	headers := DNSHeaders{
		id:                    uint16(1), // TODO : make random
		flags:                 DNSFlags{},
		questionCount:         1,
		answerRecordCount:     0,
		authorityRecordCount:  0,
		additionalRecordCount: 0,
	}

	question := DNSQuestion{
		domainName: domainName,
		type_:      1,
		class:      1,
	}

	dnsPacket, err := sendDNSQuery(dnsServerIpAddress, encodeDNSQueryAsBytes(headers, question))

	if err != nil {
		panic(errors.New(fmt.Sprintf("Error while making dns query %v", err)))
	}

	return dnsPacket
}

func getIpFromDNSPacket(packet DNSPacket) string {
	for _, answer := range packet.answers {
		if answer.type_ == TYPE_A {
			return answer.data
		}
	}
	return ""
}

func getNameServerIpFromDNSPacket(packet DNSPacket) string {
	for _, additional := range packet.additionals {
		if additional.type_ == TYPE_A {
			return additional.data
		}
	}
	return ""
}

func getNameServerDomainNameFromDNSPacket(packet DNSPacket) string {
	for _, authority := range packet.authorities {
		if authority.type_ == TYPE_NS {
			return authority.data
		}
	}
	return ""
}

func getCNameFromDNSPacket(packet DNSPacket) string {
	for _, answer := range packet.answers {
		if answer.type_ == TYPE_CNAME {
			return answer.data
		}
	}
	return ""
}

func resolveDomainNameToIp(domainName string) string {
	nameServerIP := ROOT_SERVER_IP_ADDRESS
	var resolvedIp string
	for true {
		fmt.Println(fmt.Sprintf("Query name server %s for domain name %s", nameServerIP, domainName))
		dnsPacket := sendQueryForDomainName(nameServerIP, domainName)
		// The Ip we are looking for is directly in the packet, yay !
		ip := getIpFromDNSPacket(dnsPacket)
		if ip != "" {
			resolvedIp = ip
			break
		}
		// Actually there is a cname record so we should resolve another domain name instead, let's go !
		cname := getCNameFromDNSPacket(dnsPacket)
		if cname != "" {
			resolvedIp = resolveDomainNameToIp(cname)
			break
		}
		// The server we queried does not have the response but it indicates us a server which might know the answer, directly via its IP. Great !
		nameServerIP = getNameServerIpFromDNSPacket(dnsPacket)
		if nameServerIP != "" {
			continue
		}
		// The server we queried does not have the response but it indicates us a server which might know the answer, with its domain name.
		// We can resolve that and then continue our loop to find the IP we wanted initially
		nameServerDomainName := getNameServerDomainNameFromDNSPacket(dnsPacket)
		if nameServerDomainName != "" {
			nameServerIP = resolveDomainNameToIp(nameServerDomainName)
			continue
		}

		fmt.Printf("%+v\n", dnsPacket)
		panic(errors.New("no name server domain name nor ip found in dns packet, resolving cannot continue"))
	}
	return resolvedIp
}

func main() {
	fmt.Println(resolveDomainNameToIp("www.facebook.com"))
}
