package etc

// Ported from https://github.com/stoklund/varint/blob/master/leb128.cpp .
func Leb128Encode(in []uint64) []byte {
	out := []byte{}
	_in := make([]uint64, len(in))
	copy(_in, in)

	for i := range _in {
		x := _in[i]
		for x > 127 {
			out = append(out, uint8(x|0x80))
			x >>= 7
		}
		out = append(out, uint8(x))
	}

	return out
}

func Leb128Decode(count int, in []byte) (int, []uint64) {
	x := 0
	var out []uint64
	for {
		if (count > 0) == false {
			break
		}

		count--

		_byte := in[x]
		x++

		if _byte < 128 {
			out = append(out, uint64(_byte))
			continue
		}

		var value uint64 = uint64(_byte & 0x7F)
		var shift uint64 = 7

		for _byte >= 128 {
			_byte = in[x]
			x++
			value |= uint64(_byte&0x7F) << shift
			shift += 7
		}

		out = append(out, value)
	}

	return x, out
}
