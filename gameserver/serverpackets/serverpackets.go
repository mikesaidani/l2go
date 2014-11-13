package serverpackets

import (
	"bytes"
)

func NewCryptInitPacket() []byte {

  key := []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c,	0x54, 0x87}

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x00)                       // Packet type: CruptInit
	buffer.WriteByte(0x01)                       // ?
	buffer.Write(key) // Key

	return buffer.Bytes()
}
