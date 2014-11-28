package models

import (
	"errors"
	"fmt"
	"github.com/frostwind/l2go/packets"
	"net"
)

type GameServer struct {
	Id     uint8
	Socket net.Conn
}

func NewGameServer() *GameServer {
  return &GameServer{}
}

func (g *GameServer) Receive() (opcode byte, data []byte, e error) {
	// Read the first two bytes to define the packet size
	header := make([]byte, 2)
	n, err := g.Socket.Read(header)

	if n != 2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet header.")
	}

	// Calculate the packet size
	size := 0
	size = size + int(header[0])
	size = size + int(header[1])*256

	// Allocate the appropriate size for our data (size - 2 bytes used for the length
	data = make([]byte, size-2)

	// Read the encrypted part of the packet
	n, err = g.Socket.Read(data)

	if n != size-2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet data.")
	}

	// Print the raw packet
	fmt.Printf("Raw packet : %X%X\n", header, data)

	// Extract the op code
	opcode = data[0]
	data = data[1:]
	e = nil
	return
}

func (g *GameServer) Send(data []byte) error {
	// Calculate the packet length
	length := uint16(len(data) + 2)

	// Put everything together
	buffer := packets.NewBuffer()
	buffer.WriteUInt16(length)
	buffer.Write(data)

	_, err := g.Socket.Write(buffer.Bytes())

	if err != nil {
		return errors.New("The packet couldn't be sent.")
	}

	return nil
}
