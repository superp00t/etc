package etc

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
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

func (e *Buffer) Close() error {
	return e.backend.Close()
}

func (e *Buffer) Delete() error {
	fd, ok := e.backend.(*fsBackend)
	if !ok {
		return fmt.Errorf("can't delet non-file Buffer")
	}

	e.backend.Close()

	return os.Remove(fd.path)
}

func (e *Buffer) Size() int64 {
	return e.backend.Size()
}
