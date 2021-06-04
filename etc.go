package etc

import (
	"encoding/base64"
	"fmt"
	"hash"
	"os"
)

type Buffer struct {
	flags   uint8
	backend Backend

	tmpBitsWrite    uint8
	tmpBitsWriteOfs uint8
	tmpBitsRead     uint8
	tmpBitsReadOfs  uint8
}

func NewBuffer() *Buffer {
	bf := &Buffer{
		backend: MemBackend(),
	}
	bf.init()
	return bf
}

func (b *Buffer) init() {
	b.tmpBitsWriteOfs = 8
	b.tmpBitsReadOfs = 8
}

func (b *Buffer) Erase() error {
	return b.backend.Erase()
}

func FromString(b string) *Buffer {
	return FromBytes([]byte(b))
}

// FromBytes creates a copy of the supplied slice
func FromBytes(b []byte) *Buffer {
	bf := &Buffer{}
	mb := &memBackend{
		make([]byte, len(b)),
		0,
	}
	copy(mb.buf, b)
	bf.backend = mb
	bf.init()
	return bf
}

// OfBytes creates a Buffer using the supplied slice as its backend reader
func OfBytes(b []byte) *Buffer {
	bf := &Buffer{}
	mb := &memBackend{
		b, 0,
	}
	bf.backend = mb
	bf.init()
	return bf
}

func FromBase64(s string) *Buffer {
	b, _ := base64.URLEncoding.DecodeString(s)
	return FromBytes(b)
}

func (b *Buffer) Base64() string {
	return base64.URLEncoding.EncodeToString(b.Bytes())
}

func (b *Buffer) Bytes() []byte {
	return b.backend.Bytes()
}

func (b *Buffer) Digest(hs func() hash.Hash) []byte {
	h := hs()
	h.Write(b.Bytes())
	return h.Sum(nil)
}

func (e *Buffer) Close() error {
	return e.backend.Close()
}

func (e *Buffer) Delete() error {
	fd, ok := e.backend.(*fsBackend)
	if !ok {
		return fmt.Errorf("can't delete non-file Buffer")
	}

	e.backend.Close()

	return os.Remove(fd.path)
}

func (e *Buffer) Size() int64 {
	return e.backend.Size()
}
