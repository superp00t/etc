package etc

import (
	"compress/zlib"
	"crypto/sha512"
	"encoding/base64"
	"io"
)

type Buffer struct {
	backend Backend
}

func NewBuffer() *Buffer {
	return &Buffer{
		backend: MemBackend(),
	}
}

func (b *Buffer) Flush() error {
	return b.backend.Flush()
}

func FromString(b string) *Buffer {
	return MkBuffer([]byte(b))
}

func FromBytes(b []byte) *Buffer {
	return MkBuffer(b)
}

func FromBase64(s string) *Buffer {
	b, _ := base64.URLEncoding.DecodeString(s)
	return FromBytes(b)
}

func MkBuffer(b []byte) *Buffer {
	bf := NewBuffer()
	bf.Write(b)
	return bf
}

func (b *Buffer) Base64() string {
	return base64.URLEncoding.EncodeToString(b.Bytes())
}

func (b *Buffer) Bytes() []byte {
	return b.backend.Bytes()
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

func (e *Buffer) Close() error {
	return e.backend.Close()
}
