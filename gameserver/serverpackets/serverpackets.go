package serverpackets

import (
	"bytes"
	"github.com/frostwind/l2go/packets"
)

func NewCryptInitPacket() []byte {

	key := []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x00) // Packet type: CryptInit
	buffer.WriteByte(0x01) // ?
	buffer.Write(key)      // Key

	return buffer.Bytes()
}

func NewCharListPacket() []byte {

	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x1f)                       // Packet type: CharList
  buffer.Write([]byte{0x00, 0x00, 0x00, 0x00}) // TODO

	return buffer.Bytes()
}

func NewCharTemplatePacket() []byte {

	buffer := new(packets.Buffer)
	buffer.WriteByte(0x23)   // Packet type: CharTemplate
	buffer.WriteUInt32(0x00) // We don't actually need to send the template to the client

	return buffer.Bytes()
}
