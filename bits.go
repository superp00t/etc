package etc

// Etc uses most-significant-bit first encoding of bits. For 8-byte aligned boolean values, refer to ReadBool().
func (b *Buffer) ReadMSBit() bool {
	b.tmpBitsReadOfs++
	if b.tmpBitsReadOfs > 7 {
		b.tmpBitsRead = b.ReadByte()
		b.tmpBitsReadOfs = 0
	}

	return ((b.tmpBitsRead >> (7 - b.tmpBitsReadOfs)) & 1) != 0
}

func (b *Buffer) WriteMSBit(bit bool) {
	b.tmpBitsWriteOfs--
	if bit {
		b.tmpBitsWrite |= (1 << (b.tmpBitsWriteOfs))
	}

	if b.tmpBitsWriteOfs == 0 {
		b.tmpBitsWriteOfs = 8
		b.WriteByte(b.tmpBitsWrite)
		b.tmpBitsWrite = 0
	}
}

func (b *Buffer) FlushBits() {
	b.tmpBitsRead = 0
	b.tmpBitsReadOfs = 8
	if b.tmpBitsWriteOfs < 8 {
		b.WriteByte(b.tmpBitsWrite)
	}
	b.tmpBitsWrite = 0
	b.tmpBitsWriteOfs = 8
}

func (b *Buffer) ReadLSBit() bool {
	b.tmpBitsReadOfs++
	if b.tmpBitsReadOfs > 7 {
		b.tmpBitsRead = b.ReadByte()
		b.tmpBitsReadOfs = 0
	}

	return ((b.tmpBitsRead >> (b.tmpBitsReadOfs)) & 1) != 0
}

func (b *Buffer) WriteLSBit(bit bool) {
	b.tmpBitsWriteOfs--
	if bit {
		b.tmpBitsWrite |= (1 << (7 - b.tmpBitsWriteOfs))
	}

	if b.tmpBitsWriteOfs == 0 {
		b.tmpBitsWriteOfs = 8
		b.WriteByte(b.tmpBitsWrite)
		b.tmpBitsWrite = 0
	}
}
