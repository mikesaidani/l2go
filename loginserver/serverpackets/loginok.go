package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewLoginOkPacket() []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x03)         // Packet type: LoginOk
	buffer.WriteUInt32(0x55555555) // Session id 1/2
	buffer.WriteUInt32(0x44444444) // Session id 2/2
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x000003ea)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x02)

	return buffer.Bytes()
}
