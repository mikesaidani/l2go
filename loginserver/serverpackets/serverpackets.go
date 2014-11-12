package serverpackets

import (
	"bytes"
)

func NewInitPacket() []byte {

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x00)                       // Packet type: Init
	buffer.Write([]byte{0x9c, 0x77, 0xed, 0x03}) // Session id?
	buffer.Write([]byte{0x5a, 0x78, 0x00, 0x00}) // Protocol version : 785a

	return buffer.Bytes()
}

func NewLoginOkPacket() []byte {

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x03)                       // Packet type: LoginOk
	buffer.Write([]byte{0x55, 0x55, 0x55, 0x55}) // Session id 1/2
	buffer.Write([]byte{0x44, 0x44, 0x44, 0x44}) // Session id 2/2
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	buffer.Write([]byte{0x00, 0x00, 0x03, 0xea})
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x02)

	return buffer.Bytes()
}
