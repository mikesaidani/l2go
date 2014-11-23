package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewCryptInitPacket() []byte {
	key := []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}

	buffer := packets.NewBuffer()
	buffer.WriteByte(0x00) // Packet type: CryptInit
	buffer.WriteByte(0x01) // ?
	buffer.Write(key)      // Key

	return buffer.Bytes()
}
