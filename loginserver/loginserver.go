package loginserver

import (
	"fmt"
	"github.com/frostwind/l2go/loginserver/serverpackets"
	"github.com/frostwind/l2go/packet"
	"net"
)

func handleConnection(conn net.Conn) {

	fmt.Println("A client is trying to connect...")

	fmt.Println("Building the Init packet...")

	buffer := serverpackets.NewInitPacket()
	err := packet.Send(conn, buffer, false, false)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Init packet sent.")
	}

	for {
		p, err := packet.Receive(conn)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing the connection...")
			conn.Close()
			break
		}

		switch opcode := p.GetOpcode(); opcode {
		case 00:
			buffer := serverpackets.NewLoginOkPacket()
			err := packet.Send(conn, buffer)

			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Couldn't detect the packet type.")
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
