package etc

import (
	"fmt"
	"io"
)

type dummyWriter struct {
	write io.Writer
	wpos  int64
}

// DummyWriter creates a Buffer using a Backend that can only do simple write operations.
func DummyWriter(w io.Writer) *Buffer {
	wr := &dummyWriter{
		w,
		0,
	}

	e := NewBuffer()
	e.Attach(wr)

	return e
}

func (d *dummyWriter) Write(b []byte) (int, error) {
	i, err := d.write.Write(b)
	d.wpos += int64(i)
	return i, err
}

func (d *dummyWriter) Read(b []byte) (int, error) {
	return 0, fmt.Errorf("cannot read from dummy writer")
}

func (d *dummyWriter) Erase() error {
	return fmt.Errorf("cannot erase dummy writer")
}

func (d *dummyWriter) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("cannot seek dummy writer")
}

func (d *dummyWriter) Bytes() []byte {
	return nil
}

func (d *dummyWriter) Size() int64 {
	return d.wpos
}

func (d *dummyWriter) Close() error {
	return fmt.Errorf("cannot close dummy writer")
}
