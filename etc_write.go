package etc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math"
	"math/big"
	"reflect"
	"time"
	"unicode/utf8"
)

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func (b *Buffer) WriteUString(s string) {
	b.WriteUTF8(s)
}

func (b *Buffer) WriteInvertedString(l int, s string) {
	v := reverseString(s)
	b.WriteFixedString(l, v)
}

func (b *Buffer) WriteByte(v uint8) {
	b.backend.Write([]byte{v})
}

func (b *Buffer) Write(buf []byte) (int, error) {
	return b.backend.Write(buf)
}

func (b *Buffer) WriteTypedInt(fixed bool, bits int, signed, big bool, val interface{}) {
	var v interface{}
	value := reflect.ValueOf(val)
	if value.Kind() == reflect.Ptr {
		v = value.Elem().Interface()
	} else {
		v = val
	}

	if !fixed {
		if signed {
			b.WriteInt(v.(int64))
		} else {
			b.WriteUint(v.(uint64))
		}
		return
	}

	switch bits {
	case 8:
		if signed {
			b.WriteInt8(v.(int8))
		} else {
			b.WriteByte(v.(byte))
		}
	case 16:
		if big {
			if signed {
				b.WriteBigInt16(v.(int16))
			} else {
				b.WriteBigUint16(v.(uint16))
			}
		} else {
			if signed {
				b.WriteInt16(v.(int16))
			} else {
				b.WriteUint16(v.(uint16))
			}
		}
	case 32:
		if big {
			if signed {
				b.WriteBigInt32(v.(int32))
			} else {
				b.WriteBigUint32(v.(uint32))
			}
		} else {
			if signed {
				b.WriteInt32(v.(int32))
			} else {
				b.WriteUint32(v.(uint32))
			}
		}
	case 64:
		if big {
			if signed {
				b.WriteBigInt64(v.(int64))
			} else {
				b.WriteBigUint64(v.(uint64))
			}
		} else {
			if signed {
				b.WriteInt64(v.(int64))
			} else {
				b.WriteUint64(v.(uint64))
			}
		}
	}
}

func (b *Buffer) WriteInt8(v int8) {
	b.WriteByte(uint8(v))
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

func (b *Buffer) WriteRandom(i int) *Buffer {
	by := make([]byte, i)
	_, err := io.ReadFull(rand.Reader, by)
	if err != nil {
		panic(err)
	}

	b.Write(by)
	return b
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

func (b *Buffer) WriteAt(p []byte, off int64) (int, error) {
	wpos := b.Wpos()
	b.SeekW(off)
	i, err := b.Write(p)
	b.SeekR(wpos)
	return i, err
}

func (b *Buffer) WriteBool(v bool) {
	var bit uint8
	if v {
		bit++
	}
	b.WriteByte(bit)
}

func (b *Buffer) WriteDate(t time.Time) {
	b.WriteUint(uint64(t.UnixNano() / int64(time.Millisecond)))
}

func (b *Buffer) WriteTime(t time.Time) {
	b.WriteInt(t.UnixNano())
}

func (b *Buffer) WriteRune(r rune) {
	buf := make([]byte, 8)
	n := utf8.EncodeRune(buf, r)
	b.Write(buf[:n])
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

func (b *Buffer) WriteUTF8(u string) {
	data := []byte(u)
	length := len(data)

	b.WriteUint(uint64(length))
	b.Write(data)
}
