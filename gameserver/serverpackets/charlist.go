package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewCharListPacket() []byte {
	buffer := packets.NewBuffer()
	buffer.WriteByte(0x1f)                       // Packet type: CharList
	buffer.Write([]byte{0x00, 0x00, 0x00, 0x00}) // TODO

	return buffer.Bytes()
}
