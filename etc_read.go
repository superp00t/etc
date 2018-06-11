package etc

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/unicode"
)

func (b *Buffer) Runes() []rune {
	return []rune(b.String())
}

func (b *Buffer) ReadRune() (rune, int, error) {
	if b.Available() <= 0 {
		return 0, 0, io.EOF
	}

	ahead := b.backend.Rpos()
	by := make([]byte, 8)
	b.Read(by)

	c := by[0]
	if c < utf8.RuneSelf {
		b.backend.Seek(ahead + 1)
		return rune(c), 1, nil
	}

	r, n := utf8.DecodeRune(by)
	b.backend.Seek(ahead + int64(n))

	if r == utf8.RuneError {
		return 0, 0, fmt.Errorf("Invalid utf8 sequence.")
	}

	return r, n, nil
}

func (b *Buffer) ReadInvertedString(l int) string {
	s := b.ReadFixedString(l)

	return reverseString(s)
}

func (b *Buffer) Available() int {
	return b.Len() - int(b.backend.Rpos())
}

func (b *Buffer) String() string {
	f, ok := b.backend.(*fsBackend)
	if ok {
		return fmt.Sprintf("(file @%s)", f.path)
	}

	s := string(b.Bytes())

	if utf8.ValidString(s) {
		return s
	}

	return "(non-UTF8 string)"
}

func (b *Buffer) SeekW(offset int64) {
	b.backend.SeekW(offset)
}

func (b *Buffer) SeekR(offset int64) {
	b.Seek(offset)
}

func (b *Buffer) Seek(offset int64) {
	b.backend.Seek(offset)
}

func (b *Buffer) Rpos() int64 {
	return b.backend.Rpos()
}

func (b *Buffer) Wpos() int64 {
	return b.backend.Wpos()
}

func (b *Buffer) ReadByte() uint8 {
	var bu [1]byte
	b.Read(bu[:])
	return bu[0]
}

func (b *Buffer) Len() int {
	return int(b.backend.Size())
}

func (b *Buffer) Read(buf []byte) (int, error) {
	return b.backend.Read(buf)
}

func (b *Buffer) ReadWChar(ln int) string {
	dec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	out, err := dec.Bytes(b.ReadBytes(ln))
	if err != nil {
		return ""
	}

	return string(out)
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

func (b *Buffer) ReadDate() time.Time {
	return time.Unix(0, int64(b.ReadUint())*int64(time.Millisecond))
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

func (b *Buffer) ReadBool() bool {
	bit := b.ReadByte()
	o := true
	if bit == 0 {
		o = false
	}

	return o
}

func (b *Buffer) ReadRemainder() []byte {
	// [ b e g i n | e n d ]
	//   1 2 3 4 5   6 7 9   rpos = 6

	out := make([]byte, int(b.backend.Size()-b.Rpos()+1))
	b.Read(out)
	return out
}

func (b *Buffer) ReadFixedString(i int) string {
	by := make([]byte, i)
	b.Read(by)
	return strings.TrimRight(string(by), "\x00")
}

func (b *Buffer) ReadUTF8() string {
	length := b.ReadUint()
	if length == 0 {
		return ""
	}

	e := b.ReadBytes(int(length))
	if !validateUTF8(e) {
		return ""
	}

	return string(e)
}

func (b *Buffer) ReadBoxNonce() *[24]byte {
	n := new([24]byte)
	copy(n[:], b.ReadBytes(24))
	return n
}

func (b *Buffer) ReadBoxKey() *[32]byte {
	n := new([32]byte)
	copy(n[:], b.ReadBytes(32))
	return n
}

func validateUTF8(data []byte) bool {
	e := string(data)

	return utf8.ValidString(e)
}
