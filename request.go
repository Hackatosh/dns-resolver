package main

import (
	"fmt"
	"net"
)

func sendDNSQuery(dnsServerIpAddress string, dnsQuery []byte) (DNSPacket, error) {
	//establish connection
	connection, err := net.Dial("udp", fmt.Sprintf("%s:%d", dnsServerIpAddress, 53))
	if err != nil {
		return DNSPacket{}, err
	}
	defer connection.Close()
	///send some data
	_, err = connection.Write(dnsQuery)
	if err != nil {
		return DNSPacket{}, err
	}
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return DNSPacket{}, err
	}
	return decodeBytesAsDNSPacket(buffer[:mLen]), nil
}
