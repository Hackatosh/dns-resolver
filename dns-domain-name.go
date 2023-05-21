package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

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

// This supports DNS Compression
// TODO : don’t allow loops in DNS compression (cf https://implement-dns.wizardzines.com/book/exercises.html)
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
			// The domain name is always finished after a compressed part
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
