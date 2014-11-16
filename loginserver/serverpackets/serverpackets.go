package serverpackets

import (
	"bytes"
)

const (
	REASON_ACCOUNT_IN_USE     = 0x07
	REASON_ACCESS_FAILED      = 0x04
	REASON_USER_OR_PASS_WRONG = 0x03
	REASON_SYSTEM_ERROR       = 0x01
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

func NewLoginFailPacket(reason byte) []byte {

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x01) // Packet type: LoginFail
	buffer.WriteByte(reason)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)

	return buffer.Bytes()
}

func NewServerListPacket() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x04)
	buffer.WriteByte(0x01) // Servers count
	buffer.WriteByte(0x00) // Unused

	// Server Data (Repeat for each server)
	buffer.WriteByte(0x01)                       // Server ID (Bartz)
	buffer.Write([]byte{0x7f, 0x00, 0x00, 0x01}) // Server IP address
	buffer.Write([]byte{0x61, 0x1e, 0x00, 0x00}) // Server port number
	buffer.WriteByte(0x0f)                       // Age limit
	buffer.WriteByte(0x01)                       // Is pvp allowed?
	buffer.Write([]byte{0x00, 0x00})             // How many players are online
	buffer.Write([]byte{0x10, 0x27})             // Maximum allowed players
	buffer.WriteByte(0x01)                       // Is this a testing server?

	return buffer.Bytes()
}

func NewPlayOkPacket() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x07)
	buffer.Write([]byte{0x34, 0x0b, 0x00, 0x01}) // Session Key
	buffer.Write([]byte{0x55, 0x66, 0x77, 0x88}) // Session Key 2?

	return buffer.Bytes()
}
