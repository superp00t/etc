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

func (d *dummyWriter) Flush() error {
	return fmt.Errorf("cannot flush dummy writer")
}

func (d *dummyWriter) Seek(v int64) {
	panic("cannot seek dummy writer")
}

func (d *dummyWriter) SeekW(v int64) {
	panic("cannot seek dummy writer")
}

func (d *dummyWriter) Wpos() int64 {
	return 0
}

func (d *dummyWriter) Rpos() int64 {
	panic("cannot get read position in dummy writer")
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
