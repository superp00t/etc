package etc

import (
	"fmt"
	"io"
)

type dummyReadSeeker struct {
	rdr  io.ReadSeeker
	rpos int64
}

func DummyReadSeeker(r io.ReadSeeker) *Buffer {
	dr := &dummyReadSeeker{
		r,
		0,
	}

	e := NewBuffer()
	e.Attach(dr)

	return e
}

func (d *dummyReadSeeker) Read(b []byte) (int, error) {
	i, err := d.rdr.Read(b)
	d.rpos += int64(i)
	return i, err
}

func (d *dummyReadSeeker) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("cannot write to dummy reader")
}

func (d *dummyReadSeeker) Flush() error {
	return fmt.Errorf("cannot flush dummy reader")
}

func (d *dummyReadSeeker) Seek(v int64) {
	d.rdr.Seek(v, 0)
	d.rpos = v
}

func (d *dummyReadSeeker) SeekW(v int64) {}
func (d *dummyReadSeeker) Wpos() int64 {
	return 0
}
func (d *dummyReadSeeker) Rpos() int64 {
	return d.rpos
}

func (d *dummyReadSeeker) Bytes() []byte {
	return nil
}

func (d *dummyReadSeeker) Size() int64 {
	return 0
}

func (d *dummyReadSeeker) Close() error {
	return fmt.Errorf("cannot close dummy reader")
}
