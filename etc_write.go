package etc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math"
	"math/big"
)

func (b *Buffer) WriteByte(v uint8) {
	b.backend.WriteByte(v)
}

func (b *Buffer) Write(buf []byte) (int, error) {
	for _, v := range buf {
		b.WriteByte(v)
	}
	return len(buf), nil
}

func (b *Buffer) WriteInt16(v int16) {
	b.WriteUint16(uint16(v))
}

func (b *Buffer) WriteInt32(v int32) {
	b.WriteUint32(uint32(v))
}

func (b *Buffer) WriteInt64(v int64) {
	b.WriteUint64(uint64(v))
}

/* WARNING: Big refers to Big-endian, not as in BigInteger. */
func (b *Buffer) WriteBigInt16(v int16) {
	b.WriteBigUint16(uint16(v))
}

func (b *Buffer) WriteBigInt32(v int32) {
	b.WriteBigUint32(uint32(v))
}

func (b *Buffer) WriteBigInt64(v int64) {
	b.WriteBigUint64(uint64(v))
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
	b.WriteUint(uint64(len(bf)))
	b.Write(bf)
}

func (b *Buffer) WriteFloat32(v float32) {
	b.WriteUint32(math.Float32bits(v))
}

func (b *Buffer) WriteFloat64(v float64) {
	b.WriteUint64(math.Float64bits(v))
}

func (b *Buffer) WriteRandom(i int) {
	by := make([]byte, i)
	_, err := io.ReadFull(rand.Reader, by)
	if err != nil {
		panic(err)
	}

	b.Write(by)
}

func (b *Buffer) WriteUint(v uint64) {
	b.EncodeUnsignedVarint(big.NewInt(0).SetUint64(v))
}

func (b *Buffer) WriteInt(v int64) {
	b.EncodeSignedVarint(big.NewInt(v))
}

func (b *Buffer) WriteBigInt(v *big.Int) {
	b.EncodeSignedVarint(v)
}

func (b *Buffer) WriteBool(v bool) {
	var bit uint8
	if v {
		bit++
	}
	b.WriteByte(bit)
}

func (b *Buffer) WriteFixedString(i int, v string) {
	by := make([]byte, i)
	copy(by, []byte(v))

	b.Write(by)
}

func (b *Buffer) Jump(offset int64) {
	o := b.Rpos()
	b.Seek(o + offset)
}

func (b *Buffer) Reverse(offset int64) {
	o := b.Rpos()
	b.Seek(o - offset)
}
