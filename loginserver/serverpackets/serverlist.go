package serverpackets

import (
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/packets"
	"net"
)

func NewServerListPacket(gameServers []config.GameServerType, remoteAddr string) []byte {
	buffer := new(packets.Buffer)
	buffer.WriteByte(0x04)
	buffer.WriteUInt8(uint8(len(gameServers))) // Servers count
	buffer.WriteByte(0x00)                     // Unused

	network, _, _ := net.SplitHostPort(remoteAddr)

	// Server Data (Repeat for each server)
	for index, gameserver := range gameServers {
		var ip net.IP
		if network == "127.0.0.1" {
			ip = net.ParseIP(gameserver.InternalIP).To4()
		} else {
			ip = net.ParseIP(gameserver.ExternalIP).To4()
		}

		buffer.WriteUInt8(uint8(index + 1))               // Server ID (Bartz)
		buffer.WriteByte(ip[0])                           // Server IP address 1/4
		buffer.WriteByte(ip[1])                           // Server IP address 2/4
		buffer.WriteByte(ip[2])                           // Server IP address 3/4
		buffer.WriteByte(ip[3])                           // Server IP address 4/4
		buffer.WriteUInt32(uint32(gameserver.Port))       // Server port number
		buffer.WriteByte(0x0f)                            // Age limit
		buffer.WriteByte(0x01)                            // Is pvp allowed?
		buffer.WriteUInt16(0)                             // How many players are online
		buffer.WriteUInt16(gameserver.Options.MaxPlayers) // Maximum allowed players
		if gameserver.Options.Testing == true {           // Is this a testing server?
			buffer.WriteByte(0x00)
		} else {
			buffer.WriteByte(0x01)
		}
		buffer.WriteUInt32(0x02) // Display a green clock (what is this for?)
	}

	return buffer.Bytes()
}
