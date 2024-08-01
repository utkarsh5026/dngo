package main

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer func(udpConn *net.UDPConn) {
		err := udpConn.Close()
		if err != nil {
			fmt.Println("Failed to close UDP connection:", err)
		}
	}(udpConn)

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		message, err := dns.UnMarshallMessage(buf[:size])

		if err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			continue
		}
		_, err = udpConn.WriteToUDP(message.Marshal(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
