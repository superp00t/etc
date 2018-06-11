package etc

import "io"

type memBackend struct {
	buf  []byte
	wpos int64
	rpos int64
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
	rd := 0
	for i := 0; i < len(b); i++ {
		b[i] = m.readByte()
		rd++
		if m.Rpos() >= m.Size() {
			return rd, io.EOF
		}
	}

	return rd, nil
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
