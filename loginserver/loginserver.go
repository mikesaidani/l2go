package loginserver

import (
	"bytes"
	"fmt"
	"github.com/frostwind/l2go/packet"
	"net"
)

func handleConnection(conn net.Conn) {

	fmt.Println("A client is trying to connect...")

	fmt.Println("Building the Init packet...")
	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x00)                       // Packet type: Init
	buffer.Write([]byte{0x9c, 0x77, 0xed, 0x03}) // Session id?
	buffer.Write([]byte{0x5a, 0x78, 0x00, 0x00}) // Protocol version : 785a

	err := packet.Send(conn, buffer.Bytes(), false, false)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Init packet sent.")
	}

	for {
		_, err := packet.Receive(conn)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing the connection...")
			conn.Close()
			break
		}
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
		} else {
			go handleConnection(conn)
		}

	}
}
