package main

import (
	"fmt"
)

func main() {
	dnsServer := DNSServer{
		network: "udp",
		address: "8.8.8.8",
		port:    53,
	}
	headers := DNSHeaders{
		id:                    uint16(1),
		flags:                 DNSFlags{isRecursionDesired: true},
		questionCount:         1,
		answerRecordCount:     0,
		authorityRecordCount:  0,
		additionalRecordCount: 0,
	}

	question := DNSQuestion{
		domainName: "www.example.com",
		type_:      1,
		class:      1,
	}

	response, err := makeRequest(dnsServer, encodeDNSQueryAsBytes(headers, question))

	if err != nil {
		fmt.Println(err)
		return
	}
	dnsPacket := decodeBytesAsDNSPacket(response)

	fmt.Printf("%+v\n", dnsPacket)
	fmt.Println(rawDNSRecordDataToString(dnsPacket.answers[0].data))
}
