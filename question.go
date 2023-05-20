package main

import (
	"bytes"
	"strings"
)

type DNSQuestion struct {
	domainName string
	type_      uint16 // type is a reserved keyword
	class      uint16
}

// Todo : verify that the domain name is in ascii
func encodeDomainName(domainName string) bytes.Buffer {
	encodedResult := bytes.Buffer{}
	splitted := strings.Split(domainName, ".")
	for _, part := range splitted {
		encodedResult.Write([]byte{byte(len(part))})
		encodedResult.WriteString(part)
	}
	encodedResult.Write([]byte{byte(0)})
	return encodedResult
}

func encodeDNSQuestionAsBytes(question DNSQuestion) []byte {

}
