package main

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

func buildDNSFlagsAsBytes(flags DNSFlags) uint16 {
	var flagsAsBytes uint16 = 0
	flagsAsBytes = setSingleBit(flagsAsBytes, 0, flags.isResponse)
	flagsAsBytes = setMultipleBits(flagsAsBytes, 1, uint16(flags.operationCode))
	flagsAsBytes = setSingleBit(flagsAsBytes, 5, flags.isAuthoritative)
	flagsAsBytes = setSingleBit(flagsAsBytes, 6, flags.isTruncated)
	flagsAsBytes = setSingleBit(flagsAsBytes, 7, flags.isRecursionDesired)
	flagsAsBytes = setSingleBit(flagsAsBytes, 8, flags.isRecursionAvailable)
	flagsAsBytes = setMultipleBits(flagsAsBytes, 11, uint16(flags.responseCode))
	return flagsAsBytes
}
