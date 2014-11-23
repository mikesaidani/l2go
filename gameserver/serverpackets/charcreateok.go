package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewCharCreateOkPacket() []byte {
	buffer := packets.NewBuffer()
	buffer.WriteByte(0x25)   // Packet type: CharCreateOk
	buffer.WriteUInt32(0x01) // Everything went like expected

	return buffer.Bytes()
}
