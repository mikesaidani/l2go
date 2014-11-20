package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

const (
	REASON_TOO_MANY_PLAYERS = 0x0f
	REASON_ACCESS_FAILED    = 0x04
)

func NewPlayFailPacket(reason uint32) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x06) // Packet type: PlayFail
	buffer.WriteUInt32(reason)

	return buffer.Bytes()
}
