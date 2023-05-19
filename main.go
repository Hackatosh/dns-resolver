package main

import (
	"bytes"
	"fmt"
	"strings"
)

func encode_domain_name(domain_name string) bytes.Buffer {
	encoded_result := bytes.Buffer{}
	splitted := strings.Split(domain_name, ".")
	for _, part := range splitted {
		encoded_result.Write([]byte{byte(len(part))})
		encoded_result.WriteString(part)
	}
	encoded_result.Write([]byte{byte(0)})
	return encoded_result
}

// Todo : verify that the domain name is in ascii

func main() {
	fmt.Println(encode_domain_name("google.com"))
}
