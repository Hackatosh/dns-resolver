package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type DNSQuestion struct {
	domainName string
	type_      uint16 // type is a reserved keyword
	class      uint16
}

// Todo : verify that the domain name is in ascii
func encodeDomainNameAsBytes(domainName string) []byte {
	encodedResult := bytes.Buffer{}
	splitted := strings.Split(domainName, ".")
	for _, part := range splitted {
		encodedResult.Write([]byte{byte(len(part))})
		encodedResult.WriteString(part)
	}
	encodedResult.Write([]byte{byte(0)})
	return encodedResult.Bytes()
}

func encodeDNSQuestionAsBytes(question DNSQuestion) []byte {
	bytes := encodeDomainNameAsBytes(question.domainName)
	bytes = binary.BigEndian.AppendUint16(bytes, question.type_)
	bytes = binary.BigEndian.AppendUint16(bytes, question.class)
	return bytes
}

func decodeBytesAsDomainName(data []byte, index int) (string, int) {
	parts := make([]string, 0)

	for length := data[index]; length != 0; length = data[index] {
		index += 1
		part := data[index : index+int(length)]
		index += int(length)
		parts = append(parts, string(part))
	}
	// We don't go in the loop so you need to account for the 0 length octet we just read
	index += 1

	return strings.Join(parts, "."), index
}

func decodeBytesAsQuestion(data []byte, index int) (DNSQuestion, int) {
	dnsQuestion := DNSQuestion{}
	dnsQuestion.domainName, index = decodeBytesAsDomainName(data, index)

	dnsQuestion.type_ = binary.BigEndian.Uint16(data[index : index+2])
	dnsQuestion.class = binary.BigEndian.Uint16(data[index+2 : index+4])
	return dnsQuestion, index + 4
}
