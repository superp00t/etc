package etc

type Backend interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Flush() error
	Seek(int64)
	SeekW(int64)
	Wpos() int64
	Rpos() int64
	Size() int64
	Bytes() []byte
	Close() error
}
