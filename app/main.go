package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/debug"
	"github.com/codecrafters-io/dns-server-starter-go/app/resolve"
	"net"
)

func readFromConnection(udpConn *net.UDPConn, toAddress string) {
	buf := make([]byte, 512)

	var resolver *net.Resolver
	fmt.Println("Resolver address:", toAddress)
	if toAddress != "" {
		resolver = resolve.NewResolver(toAddress)
	}

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		debug.ShowDNsPacketAsHex(buf[:size])
		message, err := resolve.HandleDnsResolution(buf[:size], resolver)

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

func main() {

	toAddress := flag.String("resolver", "", "Resolver address")
	flag.Parse()
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

	readFromConnection(udpConn, *toAddress)
}
