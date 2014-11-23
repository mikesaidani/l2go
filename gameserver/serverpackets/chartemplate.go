package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewCharTemplatePacket() []byte {
	buffer := packets.NewBuffer()
	buffer.WriteByte(0x23)   // Packet type: CharTemplate
	buffer.WriteUInt32(0x00) // We don't actually need to send the template to the client

	return buffer.Bytes()
}
