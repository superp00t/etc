package etcnet

type Token uint64

type ServerProvider interface {
	Start(address string) error // Async.
	Shutdown()
	PullData() ([]byte, Token)
	PushData(Token, []byte)
}

type ClientProvider interface {
	// Tolerance refers to the number of errors that can occur every three seconds before fails to connect
	Connect(address string, tolerance int) error
	PullData() ([]byte, error)
	PushData(data []byte) error
}
