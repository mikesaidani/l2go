package clientpackets

import (
	"github.com/frostwind/l2go/packets"
)

type Character struct {
	Name      string
	Race      uint32
	Sex       uint32
	ClassID   uint32
	STR       uint32
	CON       uint32
	DEX       uint32
	INT       uint32
	MEN       uint32
	WIT       uint32
	HairStyle uint32
	HairColor uint32
	Face      uint32
}

func NewCharacterCreate(request []byte) Character {
	var packet = packets.NewReader(request)
	var c Character

	c.Name = packet.ReadString()
	c.Race = packet.ReadUInt32()
	c.Sex = packet.ReadUInt32()
	c.ClassID = packet.ReadUInt32()
	c.INT = packet.ReadUInt32()
	c.STR = packet.ReadUInt32()
	c.CON = packet.ReadUInt32()
	c.MEN = packet.ReadUInt32()
	c.DEX = packet.ReadUInt32()
	c.WIT = packet.ReadUInt32()
	c.HairStyle = packet.ReadUInt32()
	c.HairColor = packet.ReadUInt32()
	c.Face = packet.ReadUInt32()

	return c
}
