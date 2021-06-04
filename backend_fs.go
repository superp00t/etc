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
}

func (f *fsBackend) Erase() error {
	f.file.Seek(0, 0)
	f.file.Truncate(0)
	return nil
}

func (f *fsBackend) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *fsBackend) Size() int64 {
	info, _ := f.file.Stat()
	return info.Size()
}

func (f *fsBackend) Bytes() []byte {
	orig, _ := f.Seek(0, io.SeekCurrent)
	f.Seek(0, io.SeekStart)

	out := make([]byte, int(f.Size()))
	if _, err := io.ReadFull(f.file, out); err != nil {
		panic(err)
	}

	f.Seek(orig, io.SeekStart)

	return out
}

func (f *fsBackend) Read(b []byte) (int, error) {
	return f.file.Read(b)
}

func (f *fsBackend) Write(b []byte) (int, error) {
	return f.file.Write(b)
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
