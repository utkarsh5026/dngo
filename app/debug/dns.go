package debug

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
)

func ShowDNsPacketAsHex(packet []byte) {
	fmt.Println("Header")

	offset := 0
	for i := 0; i < dns.HeaderSize; i += 2 {
		fmt.Println(BytesToHex(packet[offset : offset+2]))
		offset += 2
	}

	fmt.Println("Questions")
	fmt.Println(BytesToHex(packet[offset:]))
}
