package etc

import "io"

type dumbWriter struct {
	w    io.Writer
	wpos int64
}

func (d *dumbWriter) Rpos() int64 {
	panic("etc: cannot retrieve rpos on dummy writer")
}

func (d *dumbWriter) Wpos() int64 {
	return d.wpos
}

func (d *dumbWriter) Seek(o int64) {
	panic("etc: cannot seek on dummy writer")
}

func (d *dumbWriter) SeekW(o int64) {
	panic("etc: cannot seek on dummy writer")
}

func (d *dumbWriter) Write(b []byte) (int, error) {
	return d.w.Write(b)
}

func (d *dumbWriter) Read(b []byte) (int, error) {
	panic("etc: cannot read with dummy writer")
	return 0, io.EOF
}

func (d *dumbWriter) Bytes() []byte {
	panic("etc: cannot dump bytes with dummy writer")
}

func (d *dumbWriter) Close() error {
	panic("etc: cannot close dummy writer")
	return nil
}

func (d *dumbWriter) Size() int64 {
	return d.wpos
}

func (d *dumbWriter) Flush() error {
	panic("etc: cannot flush dummy writer")
}

func DummyWriter(w io.Writer) Backend {
	return &dumbWriter{w, 0}
}
