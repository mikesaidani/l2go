package gameserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/frostwind/l2go/gameserver/packet"
	"github.com/frostwind/l2go/gameserver/serverpackets"
	"net"
)

func read_int32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func handleConnection(conn net.Conn) {

	fmt.Println("A client is trying to connect...")

	// Init our keys
	var inputKey []byte = []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}
	var outputKey []byte = []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}

	// Receive ProtocolVersion
	p, err := packet.Receive(conn, nil)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Closing the connection...")
		conn.Close()
	}

	protocolVersion := read_int32(p.GetData())

	if protocolVersion < 419 {
		fmt.Println("Wrong protocol version ! <Min is 419>")
		conn.Close()
	}

	fmt.Println("Sending the Xor Key to the client...")

	buffer := serverpackets.NewCryptInitPacket()
	err = packet.Send(conn, buffer, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Init packet sent.")
	}

	for {
		_, err := packet.Receive(conn, inputKey)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing the connection...")
			conn.Close()
			break
		}

		switch opcode := p.GetOpcode(); opcode {
		case 0x00:
			fmt.Println("Client is requesting login to the Game Server")

			buffer := serverpackets.NewCharListPacket()
			err := packet.Send(conn, buffer, outputKey)

			if err != nil {
				fmt.Println(err)
			}

		case 0x0e:
			fmt.Println("Client is requesting character creation template")

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
