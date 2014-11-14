package packet

import (
	"fmt"
	"github.com/frostwind/l2go/loginserver/blowfish"
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

func Receive(conn net.Conn) (*packet, error) {

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

	// Decrypt the packet data using the blowfish key
	p.data, err = blowfishDecrypt(p.data, []byte("[;'.]94-31==-%&@!^+]\000"))

	if err != nil {
		return &packet{}, packetError{"An error occured while decrypting the packet data."}
	}

	// Verify our checksum...
	if check := checksum(p.data); check {
		fmt.Printf("Decrypted packet content : %X\n", p.data)
		fmt.Println("Packet checksum ok")
	} else {
		return &packet{}, packetError{"The packet checksum doesn't look right..."}
	}

	// Extract the op code
	p.opcode = p.data[0]

	return p, nil
}

func Send(conn net.Conn, data []byte, params ...bool) error {

	// Initialize our parameters
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
		// Add 4 empty bytes for the checksum
		data = append(data, []byte{0x00, 0x00, 0x00, 0x00}...)

		// Add blowfish padding
		missing := len(data) % 8

		if missing != 0 {
			for i := missing; i < 8; i++ {
				data = append(data, byte(0x00))
			}
		}

		// Finally do the checksum
		checksum(data)
	}

	if doBlowfish == true {
		var err error
		data, err = blowfishEncrypt(data, []byte("[;'.]94-31==-%&@!^+]\000"))

		if err != nil {
			return err
		}
	}

	// Add the packet length
	length := len(data) + 2
	header := []byte{byte(length) & 0xff, byte(length>>8) & 0xff}
	data = append(header, data...)

	_, err := conn.Write(data)

	if err != nil {
		return packetError{"The packet couldn't be sent."}
	}

	return nil
}

func checksum(raw []byte) bool {

	var chksum int = 0
	count := len(raw) - 8
	i := 0

	for i = 0; i < count; i += 4 {
		var ecx int = int(raw[i])
		ecx |= int(raw[i+1]) << 8
		ecx |= int(raw[i+2]) << 0x10
		ecx |= int(raw[i+3]) << 0x18
		chksum ^= ecx
	}

	var ecx int = int(raw[i])
	ecx |= int(raw[i+1]) << 8
	ecx |= int(raw[i+2]) << 0x10
	ecx |= int(raw[i+3]) << 0x18

	raw[i] = byte(chksum)
	raw[i+1] = byte(chksum >> 0x08)
	raw[i+2] = byte(chksum >> 0x10)
	raw[i+3] = byte(chksum >> 0x18)

	return ecx == chksum
}

func blowfishDecrypt(encrypted, key []byte) ([]byte, error) {

	// Initialize our cipher
	cipher, err := blowfish.NewCipher(key)

	if err != nil {
		return nil, packetError{"Couldn't initialize the blowfish cipher"}
	}

	// Check if the encrypted data is a multiple of our block size
	if len(encrypted)%8 != 0 {
		return nil, packetError{"The encrypted data is not a multiple of the block size"}
	}

	count := len(encrypted) / 8

	decrypted := make([]byte, len(encrypted))

	for i := 0; i < count; i++ {
		cipher.Decrypt(decrypted[i*8:], encrypted[i*8:])
	}

	return decrypted, nil
}

func blowfishEncrypt(decrypted, key []byte) ([]byte, error) {

	// Initialize our cipher
	cipher, err := blowfish.NewCipher(key)

	if err != nil {
		return nil, packetError{"Couldn't initialize the blowfish cipher"}
	}

	// Check if the decrypted data is a multiple of our block size
	if len(decrypted)%8 != 0 {
		return nil, packetError{"The decrypted data is not a multiple of the block size"}
	}

	count := len(decrypted) / 8

	encrypted := make([]byte, len(decrypted))

	for i := 0; i < count; i++ {
		cipher.Encrypt(encrypted[i*8:], decrypted[i*8:])
	}

	return encrypted, nil
}
