package etc

import (
	"fmt"
	"io"
)

type memBackend struct {
	buf []byte
	pos int64
}

func MemBackend() Backend {
	m := new(memBackend)
	m.buf = []byte{}
	return m
}

func (m *memBackend) Erase() error {
	m.pos = 0
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
	if int64(len(m.buf)) == m.pos {
		m.buf = append(m.buf, v)
	} else {
		m.buf[m.pos] = v
	}

	m.pos++
}

func (m *memBackend) readByte() uint8 {
	if m.pos > int64(len(m.buf))-1 {
		return 0
	}

	ch := m.buf[m.pos]
	m.pos++
	return ch
}

func (m *memBackend) Bytes() []byte {
	return m.buf
}

func (m *memBackend) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.pos = offset
		return m.pos, nil
	case io.SeekCurrent:
		if m.pos+offset > int64(len(m.buf)) {
			return 0, fmt.Errorf("etc: out of bounds buffer seek")
		}
		m.pos += offset
		return m.pos, nil
	case io.SeekEnd:
		m.pos = int64(len(m.buf))
		return m.pos, nil
	default:
		return 0, fmt.Errorf("etc: unknown seek whence")
	}
}

func (m *memBackend) Read(b []byte) (int, error) {
	ln := len(b)
	pos := int(m.pos)

	start := pos
	end := pos + ln

	var err error = nil

	if end > len(m.buf) {
		end = len(m.buf)
		err = io.EOF
	}

	bytes := copy(b, m.buf[start:end])

	m.pos = int64(end)

	return bytes, err
}

func (m *memBackend) Write(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		m.writeByte(b[i])
	}

	return len(b), nil
}

func (b *Buffer) Attach(f Backend) {
	b.backend = f
}
