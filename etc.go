package etc

import (
	"compress/zlib"
	"crypto/sha512"
	"io"
)

type Buffer struct {
	buf        []byte
	rpos, wpos int
}

func NewBuffer() *Buffer {
	return &Buffer{
		buf:  make([]byte, fragmentSize),
		rpos: 0,
		wpos: 0,
	}
}

func MkBuffer(b []byte) *Buffer {
	bf := NewBuffer()
	bf.Write(b)
	return bf
}

func (b *Buffer) Bytes() []byte {
	return b.buf[:b.wpos]
}

func (b *Buffer) Sha512Digest() []byte {
	h := sha512.New()
	h.Write(b.Bytes())
	return h.Sum(nil)
}

func ZlibCompress(input []byte) []byte {
	out := NewBuffer()

	z := zlib.NewWriter(out)
	z.Write(input)
	z.Close()
	return out.Bytes()
}

func ZlibDecompress(input []byte) ([]byte, error) {
	b := NewBuffer()
	out := NewBuffer()
	b.Write(input)

	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(out, r)
	if err != nil {
		return nil, err
	}

	r.Close()

	return out.Bytes(), nil
}
