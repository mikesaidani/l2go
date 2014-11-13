package gameserver

import (
	"bytes"
	"fmt"
	"github.com/frostwind/l2go/gameserver/packet"
	"github.com/frostwind/l2go/gameserver/serverpackets"
	"net"
)

func handleConnection(conn net.Conn) {

	fmt.Println("A client is trying to connect...")

	// Receive ProtocolVersion
	p, err := packet.Receive(conn, false)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Closing the connection...")
		conn.Close()
	}

	fmt.Printf("Protocol version : %X\n", p.GetData())

	fmt.Println("Sending the Xor Key to the client...")

	buffer := serverpackets.NewCryptInitPacket()
	err = packet.Send(conn, buffer, false)

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

		switch opcode := p.GetOpcode(); opcode {
		case 00:
			fmt.Println("Client is requesting login to the Game Server")

		default:
			fmt.Println("Couldn't detect the packet type.")
		}
	}

}

func Init() {
	ln, err := net.Listen("tcp", ":7777")
	defer ln.Close()

	if err != nil {
		fmt.Println("Couldn't initialize the Game Server")
	} else {
		fmt.Println("Game Server initialized.")
		fmt.Println("Listening on 127.0.0.1:7777.")
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
