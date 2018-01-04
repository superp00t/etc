package etc

import (
	"encoding/binary"
	"math"
)

const (
	fragmentSize = 256
)

func (b *Buffer) WriteByte(v uint8) {
	if len(b.buf) < b.wpos+1 {
		b.buf = append(b.buf, make([]byte, fragmentSize)...)
	}
	b.buf[b.wpos] = v
	b.wpos++
}

func (b *Buffer) Write(buf []byte) (int, error) {
	for _, v := range buf {
		b.WriteByte(v)
	}
	return len(buf), nil
}

func (b *Buffer) WriteUint16(v uint16) {
	d := make([]byte, 2)
	binary.LittleEndian.PutUint16(d, v)
	b.Write(d)
}

func (b *Buffer) WriteUint32(v uint32) {
	d := make([]byte, 4)
	binary.LittleEndian.PutUint32(d, v)
	b.Write(d)
}

func (b *Buffer) WriteUint64(v uint64) {
	d := make([]byte, 8)
	binary.LittleEndian.PutUint64(d, v)
	b.Write(d)
}

func (b *Buffer) WriteBigUint16(v uint16) {
	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, v)
	b.Write(d)
}

func (b *Buffer) WriteBigUint32(v uint32) {
	d := make([]byte, 4)
	binary.BigEndian.PutUint32(d, v)
	b.Write(d)
}

func (b *Buffer) WriteBigUint64(v uint64) {
	d := make([]byte, 8)
	binary.BigEndian.PutUint64(d, v)
	b.Write(d)
}

func (b *Buffer) WriteCString(v string) {
	b.Write(append([]byte(v), 0))
}

func (b *Buffer) WriteLimitedBytes(bf []byte) {
	b.WriteUint32(uint32(len(bf)))
	b.Write(bf)
}

func (b *Buffer) WriteFloat32(v float32) {
	b.WriteUint32(math.Float32bits(v))
}

func (b *Buffer) WriteFloat64(v float64) {
	b.WriteUint64(math.Float64bits(v))
}
