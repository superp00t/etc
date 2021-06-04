package etc

import (
	"fmt"
	"io"
)

type dummyReadSeeker struct {
	rdr io.ReadSeeker
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
	return i, err
}

func (d *dummyReadSeeker) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("cannot write to dummy reader")
}

func (d *dummyReadSeeker) Erase() error {
	return fmt.Errorf("cannot erase dummy reader")
}

func (d *dummyReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return d.rdr.Seek(offset, whence)
}

func (d *dummyReadSeeker) Bytes() []byte {
	return nil
}

func (d *dummyReadSeeker) Size() int64 {
	curPos, err := d.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	sizePos, err := d.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}

	d.Seek(curPos, io.SeekStart)

	return sizePos
}

func (d *dummyReadSeeker) Close() error {
	return fmt.Errorf("cannot close dummy reader")
}

func FromReadSeeker(sk io.ReadSeeker) *Buffer {
	return &Buffer{
		backend: &dummyReadSeeker{sk},
	}
}
