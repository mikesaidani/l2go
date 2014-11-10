package loginserver

import (
	"fmt"
	"github.com/frostwind/l2go/packet"
	"net"
)

func handleConnection(conn net.Conn) {

	packet_c := []byte{
		0x00, 0x9c, 0x77, 0xed,
		0x03, 0x5a, 0x78, 0x00,
		0x00}

	length := len(packet_c) + 2
	buffer := make([]byte, length)

	buffer[0] = byte(length & 0xff)
	buffer[1] = byte((length >> 8) & 0xff)
	copy(buffer[2:], packet_c)

	fmt.Println("A client is trying to connect...")
	fmt.Printf("Created an init packet[%d] = %X\n", len(buffer), buffer)

	fmt.Println("Sending the Init packet...")
	conn.Write([]byte(buffer))

	fmt.Println("Receiving the Init response")
	for {
		received := make([]byte, 65537)

		_, _ = conn.Read(received)

		fmt.Println("Decryption ..")
		_, _ = packet.Decrypt(received)
	}

}

func Init() {
	ln, err := net.Listen("tcp", ":2106")
	defer ln.Close()

	if err != nil {
		fmt.Println("Couldn't initialize the Login Server")
	} else {
		fmt.Println("Login Server initialized.")
		fmt.Println("Listening on 127.0.0.1:2106.")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Couldn't accept the incoming connection.")
			continue
		}

		go handleConnection(conn)
	}
}
