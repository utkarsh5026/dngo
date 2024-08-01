package main

import (
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		receivedHeader := dns.UnmarshalHeader(buf[:size])
		header := dns.Header{
			ID:      receivedHeader.ID,
			OpCode:  receivedHeader.OpCode,
			QR:      true,
			QDCount: 1,
			ANCount: 1,
			RD:      receivedHeader.RD,
		}

		if receivedHeader.OpCode == 0 {
			header.RCode = 0
		} else {
			header.RCode = 4
		}

		question := dns.Question{
			Name:  "codecrafters.io",
			Type:  1,
			Class: 1,
		}

		answer := dns.Answer{
			Name:  "codecrafters.io",
			Type:  1,
			Class: 1,
			TTL:   60,
		}

		message := dns.CreateDnsMessage(&header, &question, &answer)
		_, err = udpConn.WriteToUDP(message.Marshal(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
