package etc

import (
	"fmt"
	"io"
)

type dummyReader struct {
	rdr  io.ReadSeeker
	rpos int64
}

func DummyReader(r io.ReadSeeker) *Buffer {
	dr := &dummyReader{
		r,
		0,
	}

	e := NewBuffer()
	e.Attach(dr)

	return e
}

func (d *dummyReader) Read(b []byte) (int, error) {
	i, err := d.rdr.Read(b)
	d.rpos += int64(i)
	return i, err
}

func (d *dummyReader) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("cannot write to dummy reader")
}

func (d *dummyReader) Flush() error {
	return fmt.Errorf("cannot flush dummy reader")
}

func (d *dummyReader) Seek(v int64) {
	d.rdr.Seek(v, 0)
	d.rpos = v
}

func (d *dummyReader) SeekW(v int64) {}
func (d *dummyReader) Wpos() int64 {
	return 0
}
func (d *dummyReader) Rpos() int64 {
	return d.rpos
}

func (d *dummyReader) Bytes() []byte {
	return nil
}

func (d *dummyReader) Size() int64 {
	return 0
}

func (d *dummyReader) Close() error {
	return fmt.Errorf("cannot close dummy reader")
}
