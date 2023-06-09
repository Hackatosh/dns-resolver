package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type RecordType uint16

const (
	TYPE_A     RecordType = 1
	TYPE_NS    RecordType = 2
	TYPE_CNAME RecordType = 5
)

type DNSQuestion struct {
	domainName string
	type_      RecordType // type is a reserved keyword
	class      uint16
}

type DNSRecord struct {
	domainName string
	type_      RecordType // type is a reserved keyword
	class      uint16
	ttl        uint32
	data       string
}

type DNSPacket struct {
	headers     DNSHeaders
	questions   []DNSQuestion
	answers     []DNSRecord
	authorities []DNSRecord
	additionals []DNSRecord
}

func encodeDNSQuestionAsBytes(question DNSQuestion) []byte {
	domainNameAsBytes := encodeDomainNameAsBytes(question.domainName)
	domainNameAsBytes = binary.BigEndian.AppendUint16(domainNameAsBytes, uint16(question.type_))
	domainNameAsBytes = binary.BigEndian.AppendUint16(domainNameAsBytes, question.class)
	return domainNameAsBytes
}

func decodeBytesAsDNSQuestion(data []byte, index int) (DNSQuestion, int) {
	dnsQuestion := DNSQuestion{}
	dnsQuestion.domainName, index = decodeBytesAsDomainName(data, index)

	dnsQuestion.type_ = RecordType(binary.BigEndian.Uint16(data[index : index+2]))
	dnsQuestion.class = binary.BigEndian.Uint16(data[index+2 : index+4])
	return dnsQuestion, index + 4
}

// We could use net.IP for this but it would be a bit of cheating !
func rawIPToString(dnsRecordData []byte) string {
	parts := make([]string, 0)
	for _, rawByte := range dnsRecordData {
		parts = append(parts, fmt.Sprintf("%d", rawByte))
	}
	return strings.Join(parts, ".")
}

func decodeBytesAsDNSRecord(data []byte, index int) (DNSRecord, int) {
	dnsRecord := DNSRecord{}
	dnsRecord.domainName, index = decodeBytesAsDomainName(data, index)
	dnsRecord.type_ = RecordType(binary.BigEndian.Uint16(data[index : index+2]))
	dnsRecord.class = binary.BigEndian.Uint16(data[index+2 : index+4])
	dnsRecord.ttl = binary.BigEndian.Uint32(data[index+4 : index+8])
	dataLength := int(binary.BigEndian.Uint16(data[index+8 : index+10]))

	switch RecordType(dnsRecord.type_) {
	case TYPE_A:
		rawData := data[index+10 : index+10+dataLength]
		index += 10 + dataLength
		dnsRecord.data = rawIPToString(rawData)
	case TYPE_NS:
		dnsRecord.data, index = decodeBytesAsDomainName(data, index+10)
	case TYPE_CNAME:
		dnsRecord.data, index = decodeBytesAsDomainName(data, index+10)
	default:
		dnsRecord.data = "unknown dns record type so this was not decoded !!!"
		index += 10 + dataLength
	}
	return dnsRecord, index
}

func decodeBytesAsDNSPacket(data []byte) DNSPacket {
	// Init
	var index int
	dnsPacket := DNSPacket{}
	dnsPacket.questions = make([]DNSQuestion, 0)
	dnsPacket.answers = make([]DNSRecord, 0)
	dnsPacket.authorities = make([]DNSRecord, 0)
	dnsPacket.additionals = make([]DNSRecord, 0)

	// Parsing
	dnsPacket.headers, index = decodeBytesAsDNSHeaders(data)
	for i := 0; i < int(dnsPacket.headers.questionCount); i++ {
		var question DNSQuestion
		question, index = decodeBytesAsDNSQuestion(data, index)
		dnsPacket.questions = append(dnsPacket.questions, question)
	}
	for i := 0; i < int(dnsPacket.headers.answerRecordCount); i++ {
		var answer DNSRecord
		answer, index = decodeBytesAsDNSRecord(data, index)
		dnsPacket.answers = append(dnsPacket.answers, answer)
	}
	for i := 0; i < int(dnsPacket.headers.authorityRecordCount); i++ {
		var authority DNSRecord
		authority, index = decodeBytesAsDNSRecord(data, index)
		dnsPacket.authorities = append(dnsPacket.authorities, authority)
	}
	for i := 0; i < int(dnsPacket.headers.additionalRecordCount); i++ {
		var additional DNSRecord
		additional, index = decodeBytesAsDNSRecord(data, index)
		dnsPacket.additionals = append(dnsPacket.additionals, additional)
	}

	return dnsPacket
}

func encodeDNSQueryAsBytes(headers DNSHeaders, question DNSQuestion) []byte {
	encodedHeaders := encodeDNSHeadersAsBytes(headers)
	encodedQuestion := encodeDNSQuestionAsBytes(question)
	return append(encodedHeaders, encodedQuestion...)
}
