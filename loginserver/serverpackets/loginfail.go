package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewLoginFailPacket(reason uint32) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x01) // Packet type: LoginFail
	buffer.WriteUInt32(reason)

	return buffer.Bytes()
}
