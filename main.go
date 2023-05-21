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

func resolveDomainNameToIp(domainName string) string {
	nameServerIP := ROOT_SERVER_IP_ADDRESS
	var resolvedIp string
	for true {
		fmt.Println(fmt.Sprintf("Query name server %s for domain name %s", nameServerIP, domainName))
		dnsPacket := sendQueryForDomainName(nameServerIP, domainName)
		ip := getIpFromDNSPacket(dnsPacket)
		if ip != "" {
			resolvedIp = ip
			break
		}
		nameServerIP = getNameServerIpFromDNSPacket(dnsPacket)
		if nameServerIP != "" {
			continue
		}
		nameServerDomainName := getNameServerDomainNameFromDNSPacket(dnsPacket)
		if nameServerDomainName == "" {
			panic(errors.New("no name server domain name nor ip found in dns packet, resolving cannot continue"))
		}
		nameServerIP = resolveDomainNameToIp(nameServerDomainName)
	}
	return resolvedIp
}

func main() {
	fmt.Println(resolveDomainNameToIp("twitter.com"))
}
