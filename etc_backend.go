package etc

import (
	"io"
	"os"
)

func FileController(path string) (*Buffer, error) {
	b := NewBuffer()
	fb, err := FsBackend(path)
	if err != nil {
		return nil, err
	}

	b.Attach(fb)
	return b, nil
}

type Backend interface {
	ReadByte() uint8
	WriteByte(v uint8)
	Seek(offset int64)
	SeekW(offset int64)
	Wpos() int64
	Rpos() int64
	Size() int64
	Bytes() []byte
}

type memBackend struct {
	buf  []byte
	wpos int64
	rpos int64
}

type fsBackend struct {
	file *os.File
	rpos int64
	wpos int64
}

func MemBackend() Backend {
	m := new(memBackend)
	m.buf = []byte{}
	m.wpos = 0
	m.rpos = 0

	return m
}

func (m *memBackend) Size() int64 {
	return int64(len(m.buf))
}

func (m *memBackend) WriteByte(v uint8) {
	if int64(len(m.buf)) == m.Wpos() {
		m.buf = append(m.buf, v)
	} else {
		m.buf[m.wpos] = v
	}

	m.wpos++
}

func (m *memBackend) Bytes() []byte {
	return m.buf
}

func (m *memBackend) Rpos() int64 {
	return m.rpos
}

func (m *memBackend) Wpos() int64 {
	return m.wpos
}

func (m *memBackend) Seek(offset int64) {
	m.rpos = offset
}

func (m *memBackend) SeekW(offset int64) {
	m.wpos = offset
}

func (m *memBackend) ReadByte() uint8 {
	if m.Rpos() > int64(len(m.buf)) {
		return 0
	}

	ch := m.buf[m.Rpos()]
	m.rpos += 1
	return ch
}

func (f *fsBackend) SeekW(offset int64) {
	f.wpos = offset
}

func (f *fsBackend) Wpos() int64 {
	return f.wpos
}

func (f *fsBackend) Size() int64 {
	info, _ := f.file.Stat()
	return info.Size()
}

func (f *fsBackend) Rpos() int64 {
	return f.rpos
}

func (f *fsBackend) ReadByte() uint8 {
	var i [1]byte
	f.rpos += 1
	_, err := f.file.Read(i[:])
	if err != nil {
		panic(err)
	}
	return i[0]
}

func (f *fsBackend) Seek(offset int64) {
	f.rpos = offset
	f.file.Seek(offset, 0)
}

func (f *fsBackend) Bytes() []byte {
	w := f.Wpos()
	r := f.Rpos()

	f.Seek(0)

	out := make([]byte, w)
	if _, err := io.ReadFull(f.file, out); err != nil {
		panic(err)
	}

	f.Seek(r)
	f.SeekW(w)

	return out
}

func (f *fsBackend) WriteByte(v uint8) {
	f.file.Write([]byte{v})
	f.wpos += 1
}

func FsBackend(path string) (Backend, error) {
	var err error
	f := &fsBackend{}

	if _, err := os.Stat(path); err != nil {
		_, err2 := os.Create(path)
		if err2 != nil {
			return nil, err2
		}
	}

	f.file, err = os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0700)
	return f, err
}

func (b *Buffer) Attach(f Backend) {
	b.backend = f
}
