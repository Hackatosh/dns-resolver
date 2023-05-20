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
		id:                    1,
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

	encodedHeaders := encodeDNSHeadersAsBytes(headers)
	encodedQuestion := encodeDNSQuestionAsBytes(question)

	dnsQuery := append(encodedHeaders, encodedQuestion...)

	response, err := makeRequest(dnsServer, dnsQuery)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}
