package serverpackets

import (
	"github.com/frostwind/l2go/packets"
	"crypto/rand"
)

func NewLoginOkPacket() []byte {
	key := make([]byte, 8)
	_, err := rand.Read(key)

	if err != nil {
		return nil
	}

	buffer := new(packets.Buffer)
	buffer.WriteByte(0x03) // Packet type: LoginOk
	buffer.Write(key[:4])  // Session id 1/2
	buffer.Write(key[4:8]) // Session id 2/2
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x000003ea)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x00)
	buffer.WriteUInt32(0x02)

	return buffer.Bytes()
}
