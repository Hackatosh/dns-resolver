package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func decodeBytesAsDomainName(data []byte) (string, int, error) {
	parts := make([]string, 0)
	index := 12

	for length := data[index]; length != 0; length = data[index] {
		index += 1
		part := data[index : index+int(length)]
		index += int(length)
		parts = append(parts, string(part))
	}
	// We don't go in the loop so you need to account for the 0 length octet we just read
	index += 1

	return strings.Join(parts, "."), index, nil
}

func decodeBytesAsQuestion(data []byte) (DNSQuestion, error) {
	var err error
	var index int
	fmt.Println(string(data[12:]))
	dnsQuestion := DNSQuestion{}
	dnsQuestion.domainName, index, err = decodeBytesAsDomainName(data)
	if err != nil {
		return dnsQuestion, err
	}

	dnsQuestion.type_ = binary.BigEndian.Uint16(data[index : index+2])
	dnsQuestion.class = binary.BigEndian.Uint16(data[index+2 : index+4])
	return dnsQuestion, nil
}
