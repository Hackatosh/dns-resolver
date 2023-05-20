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

func hasSingleBit(n uint16, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func get4BitesNumber(n uint16, pos uint) uint16 {
	val := n & (15 << pos) // 15 is 0000 1111
	return val
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

func decodeUint16AsDNSFlags(encodedFlags uint16) DNSFlags {
	isResponse := hasSingleBit(encodedFlags, 0)
	operationCode := get4BitesNumber(encodedFlags, 1)
	isAuthoritative := hasSingleBit(encodedFlags, 5)
	isTruncated := hasSingleBit(encodedFlags, 6)
	isRecursionDesired := hasSingleBit(encodedFlags, 7)
	isRecursionAvailable := hasSingleBit(encodedFlags, 8)
	// 3 bits are reserved
	responseCode := get4BitesNumber(encodedFlags, 11)
	return DNSFlags{isResponse: isResponse, operationCode: OperationCode(operationCode), isAuthoritative: isAuthoritative, isTruncated: isTruncated, isRecursionDesired: isRecursionDesired, isRecursionAvailable: isRecursionAvailable, responseCode: ResponseCode(responseCode)}
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

func decodeBytesAsDNSHeaders(bytes []byte) DNSHeaders {
	id := binary.BigEndian.Uint16(bytes[:2])
	flags := decodeUint16AsDNSFlags(binary.BigEndian.Uint16(bytes[2:4]))
	questionCount := binary.BigEndian.Uint16(bytes[4:6])
	answerRecordCount := binary.BigEndian.Uint16(bytes[6:8])
	authorityRecordCount := binary.BigEndian.Uint16(bytes[8:10])
	additionalRecordCount := binary.BigEndian.Uint16(bytes[10:12])
	return DNSHeaders{id: id, flags: flags, questionCount: questionCount, answerRecordCount: answerRecordCount, authorityRecordCount: authorityRecordCount, additionalRecordCount: additionalRecordCount}
}
