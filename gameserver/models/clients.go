package models

import (
	"errors"
	"fmt"
	"github.com/frostwind/l2go/gameserver/crypt/xor"
	"github.com/frostwind/l2go/packets"
	"net"
)

type Client struct {
	SessionID []byte
	Socket    net.Conn
	Cipher    *xor.Cipher
}

func NewClient() *Client {
	return &Client{Cipher: xor.NewCipher()}
}

func (c *Client) Receive(params ...bool) (opcode byte, data []byte, e error) {
	doXor := true

	// Should we skip the decryption?
	if len(params) >= 1 && params[0] == false {
		doXor = false
	}

	// Read the first two bytes to define the packet size
	header := make([]byte, 2)
	n, err := c.Socket.Read(header)

	if n != 2 || err != nil {
		return 0x00, nil, errors.New(string(header))
	}

	// Calculate the packet size
	size := 0
	size = size + int(header[0])
	size = size + int(header[1])*256

	// Allocate the appropriate size for our data (size - 2 bytes used for the length
	data = make([]byte, size-2)

	// Read the encrypted part of the packet
	n, err = c.Socket.Read(data)

	if n != size-2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet data.")
	}

	// Print the raw packet
	fmt.Printf("Raw packet : %X%X\n", header, data)

	if doXor == true {
		// Decrypt the packet data using the xor key
		xor.Decrypt(data, c.Cipher.InputKey)

		// Print the decrypted packet
		fmt.Printf("Decrypted packet content : %X\n", data)

		if err != nil {
			return 0x00, nil, errors.New("An error occured while decrypting the packet data.")
		}
	}

	// Extract the opcode, and return our values
	opcode = data[0]
	data = data[1:]
	e = nil
	return
}

func (c *Client) Send(data []byte, params ...bool) error {
	doXor := true

	// Should we skip the checksum?
	if len(params) >= 1 && params[0] == false {
		doXor = false
	}

	if doXor == true {
		// Do the encryption
		xor.Encrypt(data, c.Cipher.OutputKey)
	}

	// Add the packet length
	length := uint16(len(data) + 2)
	header := packets.NewBuffer()
	header.WriteUInt16(length)

	_, err := c.Socket.Write(header.Bytes())

	if err != nil {
		return errors.New("The packet header couldn't be sent.")
	}

	_, err = c.Socket.Write(data)

	if err != nil {
		return errors.New("The packet data couldn't be sent.")
	}

	return nil
}
