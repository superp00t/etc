package etc

type Backend interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Erase() error
	Seek(offset int64, whence int) (int64, error)
	Size() int64
	Bytes() []byte
	Close() error
}
