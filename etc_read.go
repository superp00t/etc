package etc

import (
	"encoding/binary"
	"io"
	"math"
	"math/big"
	"unicode/utf8"
)

func (b *Buffer) Runes() []rune {
	return []rune(b.String())
}

func (b *Buffer) ReadRune() (rune, int, error) {
	if b.Available() <= 0 {
		return 0, 0, io.EOF
	}

	c := b.buf[b.rpos]
	if c < utf8.RuneSelf {
		b.rpos++
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(b.buf[b.wpos:])
	b.rpos += n

	return r, n, nil
}

func (b *Buffer) WriteRune(r rune) {
	buf := make([]byte, 8)
	n := utf8.EncodeRune(buf, r)
	b.Write(buf[:n])
}

func (b *Buffer) Available() int {
	return b.Len() - b.rpos
}

func (b *Buffer) String() string {
	return string(b.Bytes())
}

func (b *Buffer) SeekW(offset int) {
	b.wpos = offset
}

func (b *Buffer) SeekR(offset int) {
	b.rpos = offset
}

func (b *Buffer) Rpos() int {
	return b.rpos
}

func (b *Buffer) Wpos() int {
	return b.wpos
}

func (b *Buffer) ReadByte() uint8 {
	if b.rpos > len(b.buf)-1 {
		return 0
	}
	i := b.buf[b.rpos]
	b.rpos++
	return i
}

func (b *Buffer) Len() int {
	return b.wpos
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
func (b *Buffer) ReadInt64() int64 {
	return int64(b.ReadUint64())
}

func (b *Buffer) ReadInt32() int32 {
	return int32(b.ReadUint32())
}

func (b *Buffer) ReadInt16() int16 {
	return int16(b.ReadUint16())
}

/* WARNING: Big refers to Big-endian, not as in BigInteger. */
func (b *Buffer) ReadBigInt64() int64 {
	return int64(b.ReadBigUint64())
}

func (b *Buffer) ReadBigInt32() int32 {
	return int32(b.ReadBigUint32())
}

func (b *Buffer) ReadBigInt16() int16 {
	return int16(b.ReadBigUint16())
}

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
	return binary.BigEndian.Uint32(b.ReadBytes(4))
}

func (b *Buffer) ReadBigUint16() uint16 {
	return binary.BigEndian.Uint16(b.ReadBytes(2))
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

func (b *Buffer) ReadString(delimiter byte) string {
	var i []byte
	for {
		c := b.ReadByte()
		if c == 0 {
			break
		}

		if c == delimiter {
			break
		}
		i = append(i, c)
	}

	return string(i)
}

func (b *Buffer) ReadStringUntil(delimiter rune) string {
	var i []rune
	for {
		c, _, err := b.ReadRune()
		if err != nil {
			break
		}

		if c == delimiter {
			break
		}

		i = append(i, c)
	}

	return string(i)
}

func (b *Buffer) ReadCString() string {
	return b.ReadString(0)
}

func (b *Buffer) ReadLimitedBytes() []byte {
	ln := b.ReadUint()
	return b.ReadBytes(int(ln))
}

func (b *Buffer) ReadUint() uint64 {
	return b.DecodeUnsignedVarint(16).Uint64()
}

func (b *Buffer) ReadInt() int64 {
	return b.DecodeSignedVarint(16).Int64()
}

func (b *Buffer) ReadBigInt() *big.Int {
	return b.DecodeSignedVarint(-1)
}
