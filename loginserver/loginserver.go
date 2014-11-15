package loginserver

import (
	"fmt"
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/loginserver/packet"
	"github.com/frostwind/l2go/loginserver/serverpackets"
	"net"
)

func handleConnection(conn net.Conn) {

	fmt.Println("A client is trying to connect...")

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

		case 02:
			serverId := p.GetData()[4+4+1] // Skip the sessionId (2*4bytes) and grab the serverId

			fmt.Printf("The client wants to connect to the server : %X\n", serverId)

			buffer := serverpackets.NewPlayOkPacket()
			err := packet.Send(conn, buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 05:
			buffer := serverpackets.NewServerListPacket()
			err := packet.Send(conn, buffer)

			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Couldn't detect the packet type.")
		}
	}

}

func Init(conf config.ConfigObject) {
	ln, err := net.Listen("tcp", ":2106")
	defer ln.Close()

	if err != nil {
		fmt.Println("Couldn't initialize the Login Server")
	} else {
		fmt.Println("Login Server listening on port 2106")
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
