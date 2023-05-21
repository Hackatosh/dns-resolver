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

func decodeBytesAsDomainName(reader bytes.Reader) (string, error) {
	parts := make([]string, 0)
	var lengthAsByte byte
	var err error
	// Read next length
	lengthAsByte, err = reader.ReadByte()
	if err != nil {
		return "", err
	}
	fmt.Println(lengthAsByte)
	for uint16(lengthAsByte) != 0 {
		// Read current part and add it to the result
		part := make([]byte, uint16(lengthAsByte))
		_, err = reader.Read(part)
		if err != nil {
			return "", err
		}
		parts = append(parts, string(part))

		// Read next length
		lengthAsByte, err = reader.ReadByte()
		if err != nil {
			return "", err
		}
	}

	return strings.Join(parts, "."), nil
}

func decodeBytesAsQuestion(reader bytes.Reader) (DNSQuestion, error) {
	var err error
	dnsQuestion := DNSQuestion{}
	dnsQuestion.domainName, err = decodeBytesAsDomainName(reader)
	if err != nil {
		return dnsQuestion, err
	}
	encodedTypeAndClass := make([]byte, 4)
	_, err = reader.Read(encodedTypeAndClass)
	if err != nil {
		return dnsQuestion, err
	}

	dnsQuestion.type_ = binary.BigEndian.Uint16(encodedTypeAndClass[:2])
	dnsQuestion.class = binary.BigEndian.Uint16(encodedTypeAndClass[2:4])

	return dnsQuestion, nil
}
