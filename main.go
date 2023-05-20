package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
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

	fmt.Println(hex.EncodeToString(dnsQuery))
}
