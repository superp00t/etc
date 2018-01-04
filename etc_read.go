package etc

import (
	"encoding/binary"
	"io"
	"math"
)

func (b *Buffer) ReadByte() uint8 {
	i := b.buf[b.rpos]
	b.rpos++
	return i
}

func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Read(buf []byte) (int, error) {
	rd := 0
	for i := 0; i < len(buf); i++ {
		if b.rpos == b.Len() {
			return rd, io.EOF
		}

		buf[i] = b.ReadByte()
		rd++
	}
	return rd, nil
}

func (b *Buffer) ReadBytes(ln int) []byte {
	buf := make([]byte, ln)
	b.Read(buf)
	return buf
}

/* Integer functions */

func (b *Buffer) ReadUint64() uint64 {
	return binary.LittleEndian.Uint64(b.ReadBytes(8))
}

func (b *Buffer) ReadUint32() uint32 {
	return binary.LittleEndian.Uint32(b.ReadBytes(4))
}

func (b *Buffer) ReadUint16() uint16 {
	return binary.LittleEndian.Uint16(b.ReadBytes(2))
}

func (b *Buffer) ReadBigUint64() uint64 {
	return binary.BigEndian.Uint64(b.ReadBytes(8))
}

func (b *Buffer) ReadBigUint32() uint32 {
	return binary.LittleEndian.Uint32(b.ReadBytes(4))
}

func (b *Buffer) ReadBigUint16() uint16 {
	return binary.LittleEndian.Uint16(b.ReadBytes(2))
}

/* Floats */
func (b *Buffer) ReadFloat32() float32 {
	i := b.ReadUint32()
	bits := math.Float32frombits(i)
	return bits
}

func (b *Buffer) ReadFloat64() float64 {
	i := b.ReadUint64()
	bits := math.Float64frombits(i)
	return bits
}

func (b *Buffer) ReadCString() string {
	var i []byte
	for {
		c := b.ReadByte()
		if c == 0 {
			break
		}
		i = append(i, c)
	}

	return string(i)
}

func (b *Buffer) ReadLimitedBytes() []byte {
	ln := b.ReadUint32()
	return b.ReadBytes(int(ln))
}
