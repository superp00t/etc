package etc

import (
	"io"
	"math/big"
)

func (bf *Buffer) DecodeSignedVarint(limit int) *big.Int {
	return DecodeSignedVarint(bf, limit)
}

func DecodeSignedVarint(input io.Reader, limit int) *big.Int {
	res := big.NewInt(0)
	var more bool = true
	var shift int = 0
	var val int = 0
	off := 0
	for more {
		bt := readByte(input)
		val = int((bt) & 0x7f)

		res = clone(res).Or(res, big.NewInt(0).Lsh(big.NewInt(int64(val)), uint(shift)))
		shift += 7
		more = ((bt & 0x80) >> 7) != 0

		off++
		if limit != -1 {
			if off > limit {
				return big.NewInt(0)
			}
		}
	}

	ux := res
	nd := big.NewInt(0).And(ux, big.NewInt(0).Lsh(
		big.NewInt(1),
		uint(shift)-1))

	if nd.Cmp(big.NewInt(0)) != 0 {
		ux = clone(ux).Sub(ux, big.NewInt(int64(1<<uint(shift))))
	}

	return ux
}

func (bf *Buffer) DecodeUnsignedVarint(limit int) *big.Int {
	return DecodeUnsignedVarint(bf, limit)
}

func readByte(rd io.Reader) uint8 {
	var b [1]byte
	rd.Read(b[:])
	return b[0]
}

func DecodeUnsignedVarint(input io.Reader, limit int) *big.Int {
	_byte := readByte(input)
	if _byte < 128 {
		return big.NewInt(int64(_byte))
	}

	value := big.NewInt(0).SetUint64(uint64(_byte & 0x7F))
	var shift uint = 7

	off := 0
	for _byte >= 128 {
		off++
		if limit != -1 {
			if off > limit {
				return big.NewInt(0)
			}
		}

		_byte = readByte(input)
		value = value.Or(value, big.NewInt(0).Lsh(big.NewInt(0).SetUint64(uint64(_byte&0x7F)), shift))
		shift += 7
	}

	return value
}

func (b *Buffer) EncodeSignedVarint(x *big.Int) {
	EncodeSignedVarint(b, x)
}

func writeByte(wr io.Writer, value uint8) {
	wr.Write([]byte{value})
}

func EncodeSignedVarint(output io.Writer, x *big.Int) {
	for {
		_byte := uint8(big.NewInt(0).And(x, big.NewInt(0x7f)).Uint64())
		x = big.NewInt(0).Rsh(x, 7)

		if ((x.Cmp(big.NewInt(0)) == 0 && _byte&0x40 == 0) || (x.Cmp(big.NewInt(-1)) == 0 && _byte&0x40 != 0)) == false {
			_byte |= 0x80
		}

		writeByte(output, _byte)

		if _byte&0x80 == 0 {
			break
		}
	}
}

func (b *Buffer) EncodeUnsignedVarint(c *big.Int) {
	x := clone(c)

	i127 := big.NewInt(0).SetUint64(127)
	for x.Cmp(i127) == 1 {
		b.WriteByte(uint8(x.Uint64() | 0x80))
		x = big.NewInt(0).Rsh(x, 7)
	}

	b.WriteByte(byte(x.Uint64()))
}

func clone(x *big.Int) *big.Int {
	y := big.NewInt(0)
	y.SetBytes(x.Bytes())
	return y
}

func isZero(x *big.Int) bool {
	return x.Cmp(big.NewInt(0)) == 0
}
