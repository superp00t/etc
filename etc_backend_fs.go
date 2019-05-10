package etc

import (
	"fmt"
	"io"
	"os"
)

func FileController(path string, readOnly ...bool) (*Buffer, error) {
	b := NewBuffer()
	ro := false

	if len(readOnly) > 0 && readOnly[0] == true {
		ro = true
	}

	fb, err := FsBackend(path, ro)
	if err != nil {
		return nil, err
	}

	b.Attach(fb)
	return b, nil
}

type fsBackend struct {
	readOnly bool
	path     string
	file     *os.File
	rpos     int64
	wpos     int64
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
	if f.readOnly {
		return 0, fmt.Errorf("etc: cannot write to readonly file buffer")
	}

	f.file.Seek(f.wpos, 0)

	i, err := f.file.Write(b)
	if err == nil {
		f.wpos += int64(i)
		return i, nil
	} else {
		return i, err
	}
}

func (f *fsBackend) Close() error {
	return f.file.Close()
}

func FsBackend(path string, readOnly bool) (Backend, error) {
	var err error
	f := &fsBackend{}
	f.readOnly = readOnly
	f.path = path

	if _, err := os.Stat(path); err != nil {
		if f.readOnly == false {
			_, err2 := os.Create(path)
			if err2 != nil {
				return nil, err2
			}
		} else {
			return nil, fmt.Errorf("etc: could not open file backend: '%s'", err.Error())
		}
	}

	if f.readOnly {
		f.file, err = os.OpenFile(path, os.O_RDONLY, 0700)
	} else {
		f.file, err = os.OpenFile(path, os.O_RDWR, 0700)
	}

	return f, err
}
