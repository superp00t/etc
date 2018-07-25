package etcnet

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	Address   string
	Token     Token
	conn      *websocket.Conn
	tolerance int
}

func (w *WSClient) Connect(address string, tolerance int) error {
	w.tolerance = tolerance

	r, err := url.Parse(address)
	if err != nil {
		return err
	}

	protocol := "ws"

	if r.Query().Get("t") == "1" {
		protocol = "wss"
	}

	port := r.Port()

	w.Address = protocol + "://" + r.Host + ":" + port + "/etc"

	return w.acquireConnection()
}

func (w *WSClient) acquireConnection() error {
	t := 0
	var err error

	for {
		w.conn, err = websocket.DefaultDialer.Dial(w.Address, nil)
		if err != nil {
			if t > w.tolerance {
				return err
			}
			time.Sleep(3 * time.Second)
			t += 1
			continue
		} else {
			return nil
		}
	}
}

func (w *WSClient) PullData() ([]byte, error) {
	_, msg, err := w.conn.ReadMessage()

	if err != nil {
		err2 := w.acquireConnection()
		if err2 != nil {
			return nil, err2
		}

		return w.PullData()
	}

	return msg, nil
}

func (w *WSClient) PushData(b []byte) error {
	err := w.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		err2 := w.acquireConnection()
		if err2 != nil {
			return nil, err2
		}

		return w.PushData(b)
	}

	return msg, nil
}

func (w *WSClient) Close() error {
	if w.conn != nil {
		return w.Close()
	}

	return nil
}
