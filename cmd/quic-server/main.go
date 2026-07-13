package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	port := "9923"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("UDP echo server listening on :%s\n", port)

	buf := make([]byte, 1500)
	for {
		n, client, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		if _, err := conn.WriteToUDP(buf[:n], client); err != nil {
			continue
		}
	}
}
