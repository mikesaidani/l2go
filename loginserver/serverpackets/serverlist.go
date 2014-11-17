package serverpackets

import (
	"github.com/frostwind/l2go/packets"
)

func NewServerListPacket() []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x04)
	buffer.WriteByte(0x01) // Servers count
	buffer.WriteByte(0x00) // Unused

	// Server Data (Repeat for each server)
	buffer.WriteByte(0x01)                       // Server ID (Bartz)
	buffer.Write([]byte{0x7f, 0x00, 0x00, 0x01}) // Server IP address
	buffer.WriteUInt32(7777)                     // Server port number
	buffer.WriteByte(0x0f)                       // Age limit
	buffer.WriteByte(0x01)                       // Is pvp allowed?
	buffer.WriteUInt16(0)                        // How many players are online
	buffer.WriteUInt16(10000)                    // Maximum allowed players
	buffer.WriteByte(0x01)                       // Is this a testing server?

	return buffer.Bytes()
}
