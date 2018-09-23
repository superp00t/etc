package etcnet

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/superp00t/etc"
)

type WS_Server struct {
	Address string

	// incoming from connections
	e chan *wsEvent
	c chan error

	// sync.Map<Token, chan <-[]byte>
	oMap *sync.Map
}

type wsEvent struct {
	buf   []byte
	token Token
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (w *WS_Server) onConnect(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.BadRequest)
		return
	}

	w.handleAny(conn)
}

func (w *WS_Server) handleAny(conn *websocket.Conn, supplied Token) {
	connL := new(sync.Mutex)

	var sToken Token = supplied
	piper := make(chan struct{})

	go func() {
		select {
		case <-piper:
			outputDataI, ok := w.oMap.Load(sToken)
			if !ok {
				return
			}

			o := outputDataI.(chan []byte)

			for {
				bytes, ok := <-o
				if !ok {
					return
				}

				connL.Lock()
				err := conn.WriteMessage(websocket.BinaryMessage, bytes)
				connL.Unlock()
				if err != nil {
					return
				}
			}
		}
	}()

	for {
		var chunk []byte
		connL.Lock()
		_, message, err := conn.ReadMessage()
		if err != nil {
			time.Sleep(3 * time.Second)
			return
		}
		connL.Unlock()

		var rdr *etc.Buffer
		// Get token out of packet header.
		if len(chunk) < 16 {
			rdr = etc.FromBytes(chunk)
		} else {
			rdr = etc.FromBytes(chunk[:16])
		}

		token := Token(rdr.ReadUint())

		if sToken == 0 {
			sToken = token
			piper <- struct{}{}
		}

		go func(chunk []byte, token Token) {
			w.PullData <- &wsEvent{chunk, token}
		}(chunk, Token)
	}
}

func (w *WS_Server) init() {
	w.c = make(chan error)
	w.e = make(chan *wsEvent)
	w.oMap = new(sync.Map)
}

func (w *WS_Server) StartServer(address string) error {
	w.init()

	r := mux.NewRouter()

	r.HandleFunc("/etc", w.onConnect)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go func() {
		w.c <- http.Serve(l, r)
	}()

	select {
	case c := <-w.c:
		return c
		l.Close()
	case <-time.After(200 * time.Millisecond):
		return nil
	}
}

func (w *WS_Server) PullData() ([]byte, Token) {
	e := <-w.e
	return e.buf, e.token
}

func (w *WS_Server) SendTo(t Token, data []byte) {
	oI, ok := w.oMap.Load(t)
	if ok {
		o := oI.(chan []byte)
		o <- data
	}
}
