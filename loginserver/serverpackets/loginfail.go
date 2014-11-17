package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

const (
	REASON_ACCOUNT_IN_USE     = 0x07
	REASON_ACCESS_FAILED      = 0x04
	REASON_USER_OR_PASS_WRONG = 0x03
	REASON_SYSTEM_ERROR       = 0x01
)

func NewLoginFailPacket(reason uint32) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x01) // Packet type: LoginFail
	buffer.WriteUInt32(reason)

	return buffer.Bytes()
}
