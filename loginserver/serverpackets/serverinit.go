package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewInitPacket() []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x00)                       // Packet type: Init
	buffer.Write([]byte{0x9c, 0x77, 0xed, 0x03}) // Session id?
	buffer.Write([]byte{0x5a, 0x78, 0x00, 0x00}) // Protocol version : 785a

	return buffer.Bytes()
}
