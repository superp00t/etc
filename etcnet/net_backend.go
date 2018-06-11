package etcnet

import (
	"encoding/hex"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/superp00t/etc"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

type Event int
type Flag int

const (
	Open Event = iota
	Ping
	Message

	// Parameters
	MAX_CLOCK_DIFFERENCE = 5000
)

type Config struct {
	Flags    Flag
	Address  string
	KeyCheck func(key string) bool
}

type Agent struct {
	id                                                  []byte
	sessionPeerKey, sessionPrivateKey, sessionPublicKey *[32]byte

	peerAddr net.Addr
	ws       *websocket.Conn
	l        *Listener
	input    chan []byte
	streams  *sync.Map
}

type Listener struct {
	IP string

	c   Config
	r   *mux.Router
	udp net.PacketConn

	agents *sync.Map
	closed bool
	errc   chan error
}

func (a *Agent) transmitRaw(b []byte) error {
	if a.l != nil {
		// Retransmission
		a.l.udp.WriteTo(b, a.peerAddr)
		return nil
	}

	return a.ws.WriteMessage(websocket.BinaryMessage, b)
}

func (a *Agent) receiveRaw() ([]byte, error) {
	if a.l != nil {
		b := <-a.input
		return b, nil
	}

	_, dat, err := a.ws.ReadMessage()
	return dat, err
}

func (l *Listener) reject(a net.Addr, reason uint8) {
	b := etc.NewBuffer()
	b.WriteUint(CONN_REJECT)
	b.WriteByte(reason)
	l.udp.WriteTo(b.Bytes(), a)
}

func rejectW(ws *websocket.Conn, reason uint8) {
	b := etc.NewBuffer()
	b.WriteUint(CONN_REJECT)
	b.WriteByte(reason)
	ws.WriteMessage(websocket.BinaryMessage, b.Bytes())
}

func (l *Listener) OnNewPeer(fn func(a *Agent)) {

}

func (l *Listener) HandleWS(rw http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{} // use default options
	c, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		return
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		return
	}

	head := etc.FromBytes(message)
	op := head.ReadUint()
	if op == CONN_INIT {
		ag.handleInit(head)
	}
}

func NewListener(c Config) *Listener {
	l := new(Listener)
	l.errc = make(chan error)
	l.c = c
	l.agents = new(sync.Map)

	l.r = mux.NewRouter()
	l.r.HandleFunc("/etcnet", l.HandleWS)
	go func() {
		l.errc <- http.ListenAndServe(l.c.Address, l.r)
	}()

	go func() {
		u, err := net.ListenPacket("udp", l.c.Address)
		if err != nil {
			l.errc <- err
			return
		}

		l.udp = u

		for {
			buf := make([]byte, 65536)
			i, addr, err := l.udp.ReadFrom(buf)
			if err != nil {
				continue
			}

			head := etc.FromBytes(buf[:i])

			op := head.ReadUint()
			if op == CONN_INIT {
				ag := new(Agent)
				ag.l = l			
				ag.peerAddr = addr
				ag.streams = new(sync.Map)
				ag.input = make(chan []byte)
				
				go ag.handleInit()
				initialized := head.ReadDate()
				signkey := head.ReadBytes(32)
				sessionkey := head.ReadBoxKey()
				signature := head.ReadBytes(64)

				tnow := time.Now()

				testData := etc.NewBuffer()
				testData.WriteDate(initialized)
				testData.Write(sessionkey[:])

				if l.c.KeyCheck != nil {
					ok := l.c.KeyCheck(strings.ToUpper(hex.EncodeToString(signkey[:])))
					if !ok {
						l.reject(addr, REJECT_SIGNING)
						continue
					}
				}

				ok := ed25519.Verify(signkey, testData.Bytes(), signature)
				if !ok {
					l.reject(addr, REJECT_SIGNING)
					continue
				}

				if tnow.Sub(initialized) > ((MAX_CLOCK_DIFFERENCE) * time.Millisecond) {
					l.reject(addr, REJECT_CLOCK)
					continue
				}

			

				l.agents.Store(addr.String(), ag)
			}

			if op == CONN_DATA {
				agi, ok := l.agents.Load(addr.String())
				if !ok {
					l.reject(addr, REJECT_NO_SESSION)
					continue
				}

				ag := agi.(*Agent)
				go func() {
					ag.input <- head.ReadLimitedBytes()
				}()
			}
		}
	}()

	return l
}

func (l *Listener) Close() {
	l.closed = true
}

func (l *Listener) ListenAndServe() error {
	e := <-l.errc
	l.Close()
	return e
}
