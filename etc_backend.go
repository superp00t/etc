package etc

import (
	"fmt"
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
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Flush() error
	Seek(int64)
	SeekW(int64)
	Wpos() int64
	Rpos() int64
	Size() int64
	Bytes() []byte
	Close() error
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

func (m *memBackend) Flush() error {
	m.wpos = 0
	m.rpos = 0
	m.buf = []byte{}

	return nil
}

func (m *memBackend) Close() error {
	return nil
}

func (m *memBackend) Size() int64 {
	return int64(len(m.buf))
}

func (m *memBackend) writeByte(v uint8) {
	if int64(len(m.buf)) == m.Wpos() {
		m.buf = append(m.buf, v)
	} else {
		m.buf[m.wpos] = v
	}

	m.wpos++
}

func (m *memBackend) readByte() uint8 {
	if m.Rpos() > int64(len(m.buf))-1 {
		return 0
	}

	ch := m.buf[m.Rpos()]
	m.rpos += 1
	return ch
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

func (m *memBackend) Read(b []byte) (int, error) {
	if m.Rpos() > m.Size() {
		return 0, io.EOF
	}

	rd := 0
	for i := 0; i < len(b); i++ {
		b[i] = m.readByte()
		rd++
	}

	return rd, nil
}

func (m *memBackend) Write(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		m.writeByte(b[i])
	}

	return len(b), nil
}

func (f *fsBackend) Flush() error {
	f.file.Seek(0, 0)
	f.file.Truncate(0)
	f.wpos = 0
	f.rpos = 0
	return nil
}

func (f *fsBackend) SeekW(offset int64) {
	f.wpos = offset
	f.file.Seek(offset, 0)
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

func (f *fsBackend) Seek(offset int64) {
	f.rpos = offset
	f.file.Seek(offset, 0)
}

func (f *fsBackend) Bytes() []byte {
	w := f.Wpos()
	r := f.Rpos()

	f.Seek(0)

	out := make([]byte, int(f.Size()))
	if _, err := io.ReadFull(f.file, out); err != nil {
		panic(err)
	}

	f.Seek(r)
	f.SeekW(w)

	return out
}

func (f *fsBackend) Read(b []byte) (int, error) {
	f.file.Seek(f.rpos, 0)
	i, err := f.file.Read(b)
	if err == nil {
		f.rpos += int64(i)
	}
	return i, err
}

func (f *fsBackend) Write(b []byte) (int, error) {
	f.file.Seek(f.wpos, 0)

	i, err := f.file.Write(b)
	if err == nil {
		f.wpos += int64(i)
		return i, nil
	} else {
		fmt.Println("ERR", "=>", err)
		return i, err
	}
}

func (f *fsBackend) Close() error {
	return f.file.Close()
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

	f.file, err = os.OpenFile(path, os.O_RDWR, 0700)
	return f, err
}

func (b *Buffer) Attach(f Backend) {
	b.backend = f
}
