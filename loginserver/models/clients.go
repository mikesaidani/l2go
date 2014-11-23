package models

import (
	"crypto/rand"
	"net"
  "errors"
  "github.com/frostwind/l2go/packets"
	"github.com/frostwind/l2go/loginserver/crypt"
  "fmt"
)

type Client struct {
	Account   Account
	SessionID []byte
	Socket    net.Conn
}

func NewClient() *Client {
	id := make([]byte, 16)
	_, err := rand.Read(id)

	if err != nil {
		return nil
	}
  return &Client{SessionID: id}
}

func (c *Client) Receive() (opcode byte, data []byte, e error) {
	// Read the first two bytes to define the packet size
  header := make([]byte, 2)
	n, err := c.Socket.Read(header)

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
	n, err = c.Socket.Read(data)

	if n != size-2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet data.")
	}

	// Print the raw packet
	fmt.Printf("Raw packet : %X%X\n", header, data)

	// Decrypt the packet data using the blowfish key
	data, err = crypt.BlowfishDecrypt(data, []byte("[;'.]94-31==-%&@!^+]\000"))

	if err != nil {
		return 0x00, nil, errors.New("An error occured while decrypting the packet data.")
	}

	// Verify our checksum...
	if check := crypt.Checksum(data); check {
		fmt.Printf("Decrypted packet content : %X\n", data)
		fmt.Println("Packet checksum ok")
	} else {
		return 0x00, nil, errors.New("The packet checksum doesn't look right...")
	}

	// Extract the op code
	opcode = data[0]
  data = data[1:]
  e = nil
	return
}

func (c *Client) Send(data []byte, params ...bool) error {
	var doChecksum, doBlowfish bool = true, true

	// Should we skip the checksum?
	if len(params) >= 1 && params[0] == false {
		doChecksum = false
	}

	// Should we skip the blowfish encryption?
	if len(params) >= 2 && params[1] == false {
		doBlowfish = false
	}

	if doChecksum == true {
		// Add 4 empty bytes for the checksum new( new(
		data = append(data, []byte{0x00, 0x00, 0x00, 0x00}...)

		// Add blowfish padding
		missing := len(data) % 8

		if missing != 0 {
			for i := missing; i < 8; i++ {
				data = append(data, byte(0x00))
			}
		}

		// Finally do the checksum
		crypt.Checksum(data)
	}

	if doBlowfish == true {
		var err error
		data, err = crypt.BlowfishEncrypt(data, []byte("[;'.]94-31==-%&@!^+]\000"))

		if err != nil {
			return err
		}
	}

	// Calculate the packet length
	length := uint16(len(data) + 2)

  // Put everything together
  buffer := packets.NewBuffer()
  buffer.WriteUInt16(length)
  buffer.Write(data)

	_, err := c.Socket.Write(buffer.Bytes())

	if err != nil {
		return errors.New("The packet couldn't be sent.")
	}

	return nil
}
