package main

import (
	"fmt"
	"net"
)

type DNSServer struct {
	network string
	address string
	port    uint16
}

func makeRequest(dnsServer DNSServer, requestInBytes []byte) ([]byte, error) {
	//establish connection
	connection, err := net.Dial(dnsServer.network, fmt.Sprintf("%s:%d", dnsServer.address, dnsServer.port))
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	///send some data
	_, err = connection.Write(requestInBytes)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 2048)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:mLen], nil
}
