package packet

import (
	"fmt"
	"net"
)

// A packet is composed of:
// - A 2 bytes header (defining the size)
// - An opcode describing the packet type
// - Packet data
type packet struct {
	header []byte
	opcode byte
	data   []byte
}

type packetError struct {
	message string
}

func (e packetError) Error() string {
	return fmt.Sprintf("%v", e.message)
}

func (p *packet) GetOpcode() byte {
	return p.opcode
}

func (p *packet) GetData() []byte {
	return p.data[1:]
}

func Receive(conn net.Conn, params ...bool) (*packet, error) {

	var key []byte = []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}

	// Initialize our parameters
	var doXor bool = true

	// Should we skip the checksum?
	if len(params) == 1 && params[0] == false {
		doXor = false
	}

	// Init our packet struct
	p := new(packet)

	// Read the first two bytes to define the packet size
	p.header = make([]byte, 2)
	n, err := conn.Read(p.header)

	if n != 2 || err != nil {
		return &packet{}, packetError{"An error occured while reading the packet header."}
	}

	// Calculate the packet size
	size := 0
	size = size + int(p.header[0])
	size = size + int(p.header[1])*256

	// Allocate the appropriate size for our data (size - 2 bytes used for the length
	p.data = make([]byte, size-2)

	// Read the encrypted part of the packet
	n, err = conn.Read(p.data)

	if n != size-2 || err != nil {
		return &packet{}, packetError{"An error occured while reading the packet data."}
	}

	// Print the raw packet
	fmt.Printf("Raw packet : %X%X\n", p.header, p.data)

	if doXor == true {
		// Decrypt the packet data using the blowfish key
		xorDecrypt(p.data, key)

		// Print the decrypted packet
		fmt.Printf("Decrypted packet content : %X\n", p.data)

		if err != nil {
			return &packet{}, packetError{"An error occured while decrypting the packet data."}
		}
	}

	// Extract the op code
	p.opcode = p.data[0]

	return p, nil
}

func Send(conn net.Conn, data []byte, params ...bool) error {

	var key []byte = []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}

	// Initialize our parameters
	var doXor bool = true

	// Should we skip the checksum?
	if len(params) == 1 && params[0] == false {
		doXor = false
	}

	if doXor == true {
		// Finally do the checksum
		xorEncrypt(data, key)
	}

	// Add the packet length
	length := len(data) + 2
	header := []byte{byte(length) & 0xff, byte(length>>8) & 0xff}
	//data = append(header, data...)

	_, err := conn.Write(header)

	if err != nil {
		return packetError{"The packet header couldn't be sent."}
	}

	_, err = conn.Write(data)

	if err != nil {
		return packetError{"The packet data couldn't be sent."}
	}

	return nil
}

func xorDecrypt(raw, key []byte) {

	temp := 0
	j := 0
	length := len(raw)

	for i := 0; i < length; i++ {
		temp2 := int(raw[i])
		raw[i] = byte(temp2) ^ (key[j] & 0xff) ^ byte(temp)
		j = j + 1
		temp = temp2

		if j > 7 {
			j = 0
		}
	}

	var old int = int(key[0])
	old |= int(key[1]) << 8
	old |= int(key[2]) << 0x10
	old |= int(key[3]) << 0x18

	old += len(raw)

	key[0] = byte(old)
	key[1] = byte(old >> 0x08)
	key[2] = byte(old >> 0x10)
	key[3] = byte(old >> 0x18)
}

func xorEncrypt(raw, key []byte) {

	temp := 0
	j := 0
	length := len(raw)

	for i := 0; i < length; i++ {
		temp2 := int(raw[i])
		raw[i] = byte(temp2) ^ (key[j] & 0xff) ^ byte(temp)
		j = j + 1
		temp = int(raw[i])

		if j > 7 {
			j = 0
		}
	}

	var old int = int(key[0])
	old |= int(key[1]) << 8
	old |= int(key[2]) << 0x10
	old |= int(key[3]) << 0x18

	old += len(raw)

	key[0] = byte(old)
	key[1] = byte(old >> 0x08)
	key[2] = byte(old >> 0x10)
	key[3] = byte(old >> 0x18)
}
