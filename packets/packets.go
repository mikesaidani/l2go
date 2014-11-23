package packets

import (
	"bytes"
	"encoding/binary"
)

type Buffer struct {
	bytes.Buffer
}

func NewBuffer() *Buffer {
  return &Buffer{}
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

type Reader struct {
	*bytes.Reader
}

func NewReader(buffer []byte) *Reader {
	return &Reader{bytes.NewReader(buffer)}
}

func (r *Reader) ReadBytes(number int) []byte {
	buffer := make([]byte, number)
	n, _ := r.Read(buffer)
	if n < number {
		return []byte{}
	}

	return buffer
}

func (r *Reader) ReadUInt64() uint64 {
	var result uint64

	buffer := make([]byte, 8)
	n, _ := r.Read(buffer)
	if n < 8 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadUInt32() uint32 {
	var result uint32

	buffer := make([]byte, 4)
	n, _ := r.Read(buffer)
	if n < 4 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadUInt16() uint16 {
	var result uint16

	buffer := make([]byte, 2)
	n, _ := r.Read(buffer)
	if n < 2 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadUInt8() uint8 {
	var result uint8

	buffer := make([]byte, 1)
	n, _ := r.Read(buffer)
	if n < 1 {
		return 0
	}

	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.LittleEndian, &result)

	return result
}

func (r *Reader) ReadString() string {
	var result []byte
  var first_byte, second_byte byte

  for {
    first_byte, _ = r.ReadByte()
    second_byte, _ = r.ReadByte()
    if first_byte == 0x00 && second_byte == 0x00 {
      break
    } else {
      result = append(result, first_byte, second_byte)
    }
  }

	return string(result)
}
