package etc

import (
	"fmt"
	"io"
)

type dummyReader struct {
	rdr  io.Reader
	rpos int64
	size int64
}

func DummyReader(r io.Reader, size int64) *Buffer {
	dr := &dummyReader{
		r,
		0,
		size,
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

func (d *dummyReader) Erase() error {
	return fmt.Errorf("cannot erase dummy reader")
}

func (d *dummyReader) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("can't seek in dummy reader")
}

func (d *dummyReader) Bytes() []byte {
	return nil
}

func (d *dummyReader) Size() int64 {
	return d.size
}

func (d *dummyReader) Close() error {
	return fmt.Errorf("cannot close dummy reader")
}
