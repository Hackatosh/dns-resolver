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

type DNSRecord struct {
	domainName string
	type_      uint16 // type is a reserved keyword
	class      uint16
	ttl        uint32
	data       []byte
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
	length := data[index]
	index += 1
	for length != 0 {
		// data is compressed
		if length&0xc0 != 0 {
			pointer := int(binary.BigEndian.Uint16(data[index-1:index+1]) & 0x3fff)
			part, _ := decodeBytesAsDomainName(data, pointer)
			parts = append(parts, part)
			// Pointer is 2 bytes length
			index += 1
			break
		} else {
			part := string(data[index : index+int(length)])
			parts = append(parts, part)
			index += int(length)
			length = data[index]
			index += 1
		}

	}

	return strings.Join(parts, "."), index
}

func decodeBytesAsDNSQuestion(data []byte, index int) (DNSQuestion, int) {
	dnsQuestion := DNSQuestion{}
	dnsQuestion.domainName, index = decodeBytesAsDomainName(data, index)

	dnsQuestion.type_ = binary.BigEndian.Uint16(data[index : index+2])
	dnsQuestion.class = binary.BigEndian.Uint16(data[index+2 : index+4])
	return dnsQuestion, index + 4
}

func decodeBytesAsDNSRecord(data []byte, index int) (DNSRecord, int) {
	dnsRecord := DNSRecord{}
	dnsRecord.domainName, index = decodeBytesAsDomainName(data, index)
	dnsRecord.type_ = binary.BigEndian.Uint16(data[index : index+2])
	dnsRecord.class = binary.BigEndian.Uint16(data[index+2 : index+4])
	dnsRecord.ttl = binary.BigEndian.Uint32(data[index+4 : index+8])
	dataLength := int(binary.BigEndian.Uint16(data[index+8 : index+10]))
	dnsRecord.data = data[index+10 : index+10+dataLength]
	index += 10 + dataLength
	return dnsRecord, index
}
