package main

import (
	"bytes"
	"fmt"
	"net"
)

type DNSServer struct {
	network string
	address string
	port    uint16
}

func makeRequest(dnsServer DNSServer, requestInBytes []byte) (bytes.Reader, error) {
	//establish connection
	connection, err := net.Dial(dnsServer.network, fmt.Sprintf("%s:%d", dnsServer.address, dnsServer.port))
	if err != nil {
		return bytes.Reader{}, err
	}
	defer connection.Close()
	///send some data
	_, err = connection.Write(requestInBytes)
	if err != nil {
		return bytes.Reader{}, err
	}
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return bytes.Reader{}, err
	}
	return *bytes.NewReader(buffer[:mLen]), nil
}
