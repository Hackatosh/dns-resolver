package main

import (
	"encoding/binary"
)

type OperationCode uint16

const (
	QUERY OperationCode = iota
	IQUERY
	STATUS
	reserved
	NOTIFY
	UPDATE
)

type ResponseCode uint16

const (
	NO_ERROR ResponseCode = iota
	FORMAT_ERROR
	SERVER_FAILURE
	NAME_ERROR
	NOT_IMPLEMENTED
	REFUSED
	YX_DOMAIN
	YX_RR_SET
	NX_RR_SET
	NOT_AUTH
	NOT_ZONE
)

type DNSFlags struct {
	isResponse           bool
	operationCode        OperationCode
	isAuthoritative      bool
	isTruncated          bool
	isRecursionDesired   bool
	isRecursionAvailable bool
	responseCode         ResponseCode
}

type DNSHeaders struct {
	id                    uint16
	flags                 DNSFlags
	questionCount         uint16
	answerRecordCount     uint16
	authorityRecordCount  uint16
	additionalRecordCount uint16
}

func setSingleBit(flagsAsBytes uint16, pos uint, bitSet bool) uint16 {
	var bitSetVar uint16
	if bitSet {
		bitSetVar = 1
	}
	flagsAsBytes |= bitSetVar << pos
	return flagsAsBytes
}

func setMultipleBits(flagsAsBytes uint16, pos uint, bitesToSet uint16) uint16 {
	flagsAsBytes |= (bitesToSet << pos)
	return flagsAsBytes
}

func encodeDNSFlagsAsUint16(flags DNSFlags) uint16 {
	var flagsAsBytes uint16 = 0
	flagsAsBytes = setSingleBit(flagsAsBytes, 0, flags.isResponse)
	flagsAsBytes = setMultipleBits(flagsAsBytes, 1, uint16(flags.operationCode))
	flagsAsBytes = setSingleBit(flagsAsBytes, 5, flags.isAuthoritative)
	flagsAsBytes = setSingleBit(flagsAsBytes, 6, flags.isTruncated)
	flagsAsBytes = setSingleBit(flagsAsBytes, 7, flags.isRecursionDesired)
	flagsAsBytes = setSingleBit(flagsAsBytes, 8, flags.isRecursionAvailable)
	// 3 bits are reserved
	flagsAsBytes = setMultipleBits(flagsAsBytes, 11, uint16(flags.responseCode))
	return flagsAsBytes
}

func decodeBytesAsDNSFlags(encodedFlags []byte) DNSFlags {
	dnsFlags := DNSFlags{}
	// Byte are in reverse order regarding the spec !
	// Probably some big endian / little endian shenanigans
	// This has helped : https://github.com/google/gopacket/blob/master/layers/dns.go
	dnsFlags.isResponse = encodedFlags[0]&0x80 != 0
	dnsFlags.operationCode = OperationCode((encodedFlags[0] >> 3) & 0x0F)
	dnsFlags.isAuthoritative = encodedFlags[0]&0x04 != 0
	dnsFlags.isTruncated = encodedFlags[0]&0x02 != 0
	dnsFlags.isRecursionDesired = encodedFlags[0]&0x01 != 0
	dnsFlags.isRecursionAvailable = encodedFlags[1]&0x80 != 0
	// 3 bits are reserved
	dnsFlags.responseCode = ResponseCode(encodedFlags[1] & 0xF)
	return dnsFlags
}

func encodeDNSHeadersAsBytes(headers DNSHeaders) []byte {
	bytes := make([]byte, 0)
	bytes = binary.BigEndian.AppendUint16(bytes, headers.id)
	bytes = binary.BigEndian.AppendUint16(bytes, encodeDNSFlagsAsUint16(headers.flags))
	bytes = binary.BigEndian.AppendUint16(bytes, headers.questionCount)
	bytes = binary.BigEndian.AppendUint16(bytes, headers.answerRecordCount)
	bytes = binary.BigEndian.AppendUint16(bytes, headers.authorityRecordCount)
	bytes = binary.BigEndian.AppendUint16(bytes, headers.additionalRecordCount)
	return bytes
}

func decodeBytesAsDNSHeaders(data []byte) (DNSHeaders, int) {
	dnsHeaders := DNSHeaders{}
	encodedDNSHeaders := data[:12]
	dnsHeaders.id = binary.BigEndian.Uint16(encodedDNSHeaders[:2])
	dnsHeaders.flags = decodeBytesAsDNSFlags(encodedDNSHeaders[2:4])
	dnsHeaders.questionCount = binary.BigEndian.Uint16(encodedDNSHeaders[4:6])
	dnsHeaders.answerRecordCount = binary.BigEndian.Uint16(encodedDNSHeaders[6:8])
	dnsHeaders.authorityRecordCount = binary.BigEndian.Uint16(encodedDNSHeaders[8:10])
	dnsHeaders.additionalRecordCount = binary.BigEndian.Uint16(encodedDNSHeaders[10:12])
	return dnsHeaders, 12
}
