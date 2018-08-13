package etcnet

import "errors"

const (
	WS  = 1
	UDP = 2
)

type CCreator func() ClientProvider

var (
	TorAddress = "127.0.0.1:9050"

	cmethods = map[int]CCreator{
		WS: CCreator {
			return new(WSClient)
		},
		UDP: CCreator {
			return new(UDPClient)
		}
	}

	errNoClientProvider = errors.New("etcnet: invalid client provider")
)

type Opts struct {
	Address    string
	ClientType int
	Tolerance  int
}

type Client struct {
	Opts *Opts
}

type Conn struct {
	client bool
	Prov ClientProvider
	c    *Client
}

func Dial(address string, o *Opts) (*Conn, error) {
	if o == nil {
		o = &Opts{
			Address:    address,
			ClientType: UDP,
			Tolerance:  10,
		}
	}

	o.Address = address

	c := &Client{o}

	return c.Dial()
}

func (c *Client) Dial() (*Conn, error) {
	conn := new(Conn)
	conn.c = c
	conn.Prov = cmethods[c.Opts.ClientType]

	if conn.Prov == nil {
		return nil, errNoClientProvider
	}

	err := conn.Prov.Connect(c.Opts.Address, c.Opts.Tolerance)
	if err != nil {
		return nil, err
	}

	conn.Prov.PushData(ControlPacket {
		MessageID: 0,
		Flags:     CONTROL_ACQUIRE_TOKEN,
	})
}
