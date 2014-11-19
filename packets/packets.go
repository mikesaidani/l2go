package packets

import (
	"bytes"
	"encoding/binary"
)

type Buffer struct {
	bytes.Buffer
}

func (b *Buffer) WriteUInt64(value uint64) {
	binary.Write(b, binary.LittleEndian, value)
}

func (b *Buffer) WriteUInt32(value uint32) {
	binary.Write(b, binary.LittleEndian, value)
}

func (b *Buffer) WriteUInt16(value uint16) {
	binary.Write(b, binary.LittleEndian, value)
}

func (b *Buffer) WriteUInt8(value uint8) {
	binary.Write(b, binary.LittleEndian, value)
}

func (b *Buffer) WriteFloat64(value float64) {
	binary.Write(b, binary.LittleEndian, value)
}

func (b *Buffer) WriteFloat32(value float32) {
	binary.Write(b, binary.LittleEndian, value)
}
