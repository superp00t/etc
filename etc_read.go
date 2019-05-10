package etc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

func (b *Buffer) Runes() []rune {
	return []rune(b.String())
}

func (b *Buffer) ReadRune() (rune, int, error) {
	const (
		// first byte of a 2-byte encoding starts 110 and carries 5 bits of data
		b2Lead = 0xC0 // 1100 0000
		b2Mask = 0x1F // 0001 1111

		// first byte of a 3-byte encoding starts 1110 and carries 4 bits of data
		b3Lead = 0xE0 // 1110 0000
		b3Mask = 0x0F // 0000 1111

		// first byte of a 4-byte encoding starts 11110 and carries 3 bits of data
		b4Lead = 0xF0 // 1111 0000
		b4Mask = 0x07 // 0000 0111

		// non-first bytes start 10 and carry 6 bits of data
		mbLead = 0x80 // 1000 0000
		mbMask = 0x3F // 0011 1111
	)

	if b.Available() <= 0 {
		return 0, 0, io.EOF
	}

	header := b.ReadByte()
	if header < mbLead {
		return rune(header), 1, nil
	} else if header < b3Lead {
		b1 := b.ReadByte()
		return rune(header&b2Mask)<<6 | rune(b1&mbMask), 2, nil
	} else if header < b4Lead {
		b1 := b.ReadByte()
		b2 := b.ReadByte()
		return rune(header&b3Mask)<<12 |
			rune(b1&mbMask)<<6 |
			rune(b2&mbMask), 3, nil
	} else {
		b1 := b.ReadByte()
		b2 := b.ReadByte()
		b3 := b.ReadByte()
		return rune(header&b4Mask)<<18 |
			rune(b1&mbMask)<<12 |
			rune(b2&mbMask)<<6 |
			rune(b3&mbMask), 4, nil
	}
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
	char := make([]uint16, ln)

	for i := 0; i < ln; i++ {
		char[i] = b.ReadUint16()
	}

	return string(utf16.Decode(char))
}

func (b *Buffer) ReadBytes(ln int) []byte {
	buf := make([]byte, ln)
	b.Read(buf)
	return buf
}

/* Integer functions */
func (b *Buffer) ReadTypedInt(fixed bool, bits int, signed, big bool, v interface{}) {
	var val interface{}

	if !fixed {
		if signed {
			val = b.ReadInt()
		} else {
			val = b.ReadUint()
		}
	} else {

		switch bits {
		case 8:
			if signed {
				val = b.ReadInt8()
			} else {
				val = b.ReadByte()
			}
		case 16:
			if big {
				if signed {
					val = b.ReadBigInt16()
				} else {
					val = b.ReadBigUint16()
				}
			} else {
				if signed {
					val = b.ReadInt16()
				} else {
					val = b.ReadUint16()
				}
			}
		case 32:
			if big {
				if signed {
					val = b.ReadBigInt32()
				} else {
					val = b.ReadBigUint32()
				}
			} else {
				if signed {
					val = b.ReadInt32()
				} else {
					val = b.ReadUint32()
				}
			}
		case 64:
			if big {
				if signed {
					val = b.ReadBigInt64()
				} else {
					val = b.ReadBigUint64()
				}
			} else {
				if signed {
					val = b.ReadInt64()
				} else {
					val = b.ReadUint64()
				}
			}
		}
	}

	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(val))
	// fmt.Println("Set data to", reflect.ValueOf(v).Elem().Interface())

	// fmt.Println("Data =>", reflect.ValueOf(val).String())

	// fmt.Println("Setting value of ", reflect.ValueOf(v).Elem().String())
	// reflect.ValueOf(v).Elem().Set(reflect.ValueOf(val))
	// fmt.Println("Set value of ", reflect.ValueOf(v).Elem().String())
}

func (b *Buffer) ReadInt8() int8 {
	return int8(b.ReadByte())
}

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

func (b *Buffer) ReadAt(p []byte, off int64) (int, error) {
	rpos := b.Rpos()
	b.SeekR(off)
	i, err := b.Read(p)
	b.SeekR(rpos)
	return i, err
}

func (b *Buffer) Find(str []byte) (int64, error) {
	const bufSize = 4096 * 2

	var dat [bufSize]byte
	var scanBuf = make([]byte, len(str))

	for {
		pos := b.Rpos()

		i, err := b.Read(dat[:])
		if err != nil && err != io.EOF {
			return 0, err
		}

		npos := b.Rpos()

		for x := 0; x < i; x++ {
			u := dat[x]

			if u == str[0] {
				b.SeekR(pos + int64(x))
				b.Read(scanBuf)
				if bytes.Equal(scanBuf, str) {
					return pos + int64(x), nil
				} else {
					b.SeekR(npos)
				}
			}
		}

		if err == io.EOF {
			return 0, err
		}

		for x := 0; x < bufSize; x++ {
			dat[x] = 0
		}

		for x := 0; x < len(str); x++ {
			scanBuf[x] = 0
		}
	}
}

func (b *Buffer) ReadUntilToken(s string) (string, error) {
	tmp := []rune{}
	str := []rune(s)

	offset := 0

	for {
		if offset > len(str) {
			return "", fmt.Errorf("etc: offset error")
		}

		rn, _, err := b.ReadRune()
		if err != nil {
			return "", err
		}

		if rn == str[0] {
			f := false
			ct := []rune{}
			for x := 1; ; x++ {
				r, _, err := b.ReadRune()
				if err != nil {
					return "", err
				}

				ct = append(ct, r)

				if x > len(str)-1 {
					break
				}

				if str[x] != r {
					f = true
					break
				}
			}

			if f {
				tmp = append(tmp, rn)
				tmp = append(tmp, ct...)
				continue
			} else {
				b.SeekR(b.Rpos() - 1)
				return string(tmp), nil
			}
		} else {
			tmp = append(tmp, rn)
		}
	}
}

func (b *Buffer) ReadString(delimiter byte) (string, error) {
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

	return string(i), nil
}

func (b *Buffer) ReadUString() string {
	str := b.ReadLimitedBytes()

	if !validateUTF8(str) {
		return ""
	}

	return string(str)
}

// ReadDate returns a timestamp based on
func (b *Buffer) ReadDate() time.Time {
	return time.Unix(0, int64(b.ReadUint())*int64(time.Millisecond))
}

func (b *Buffer) ReadTime() time.Time {
	return time.Unix(0, b.ReadInt())
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
	s, _ := b.ReadString(0)
	return s
}

func (b *Buffer) ReadLimitedBytes() []byte {
	ln := b.ReadUint()
	if uint64(b.Available()) < ln {
		return []byte{}
	}

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

func (b *Buffer) ToString() string {
	return string(b.Bytes())
}

func runeSize(r rune) int {
	_, sz := utf8.DecodeRuneInString(string([]rune{r}))
	return sz
}
