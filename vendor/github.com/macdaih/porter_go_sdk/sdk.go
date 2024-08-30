package portergosdk

import (
	"context"
	"fmt"
	"net"
	"time"
)

type credential struct {
	authMethod string
	usr        *string
	pwd        *string
}

type SubscribeCallback func() error

type PorterClient struct {
	serverHost string

	clientID string

	keepAlive uint16

	willFlag uint8

	cleanStart bool
	//will       bool
	//retain     bool
	//pwdFlag bool
	//usrFlag bool
	qos uint8

	nextPacketID uint16

	creds *credential
	conn  *net.TCPConn

	messageHandler func([]byte) error

	endState chan struct{}
}

type Option func(c *PorterClient)

func WithID(id string) Option {
	return func(c *PorterClient) {
		c.clientID = id
	}
}

func WithBasicCredentials(user string, pwd string) Option {
	return func(c *PorterClient) {
		c.creds = &credential{
			authMethod: PasswordMethod,
			usr:        &user,
			pwd:        &pwd,
		}
	}
}

func WithCallBack(fn func(b []byte) error) Option {
	return func(c *PorterClient) {
		c.messageHandler = fn
	}
}

func NewClient(
	serverHost string,
	keepAlive uint16,
	options ...Option,
) *PorterClient {
	es := make(chan struct{}, 1)

	pc := PorterClient{
		serverHost:     serverHost,
		keepAlive:      keepAlive,
		endState:       es,
		messageHandler: func(_ []byte) error { return nil },
	}

	for _, fn := range options {
		fn(&pc)
	}

	return &pc
}

func (pc *PorterClient) connect(ctx context.Context) error {

	msg, err := buildConnect(
		pc.clientID,
		pc.keepAlive,
		pc.creds,
	)

	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp4", pc.serverHost)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	pc.conn = conn

	// no closed conn
	if _, err := pc.conn.Write(msg); err != nil {
		return err
	}

	received := make([]byte, 1024)
	if _, err = pc.conn.Read(received); err != nil {
		return err
	}

	// handle connack
	res := received[0]
	if res != 0x20 {
		return fmt.Errorf("failed connack response : received %s", parseCode(res))
	}

	ka := time.Duration(pc.keepAlive) * time.Second

	go func() {
		tmark := time.Now().Add(ka)
		for {

			buff := make([]byte, 1024)
			if _, err := pc.conn.Read(buff); err != nil {
				// TODO use err in struct end state
				fmt.Println(err)
				pc.endState <- struct{}{}
				return
			}

			if err := pc.readMessage(buff); err != nil {
				fmt.Println(err)
				pc.endState <- struct{}{}
				return
			}

			select {
			case <-ctx.Done():
				pc.endState <- struct{}{}
				return
			default:
				if n := time.Now(); n.After(tmark) {
					tmark = n.Add(ka)

					ping := []byte{0xC0, 0}
					if _, err := pc.conn.Write(ping); err != nil {
						pc.endState <- struct{}{}
						return
					}
				}
			}
		}
	}()

	return nil
}

func (pc *PorterClient) Subscribe(ctx context.Context, topics []string) error {
	if err := pc.connect(ctx); err != nil {
		return err
	}

	msg, err := buildSubscribe(topics, 1)
	if err != nil {
		return err
	}
	//

	if pc.conn == nil {
		return fmt.Errorf("failed to perform subscription : client disconnected")
	}

	if _, err := pc.conn.Write(msg); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		fmt.Println("ctx done")
	case <-pc.endState:
	}

	fmt.Println("client flow done")
	pc.conn.Close()

	return nil
}

func (pc *PorterClient) readMessage(pkt []byte) error {

	switch pkt[0] {
	case 0xe0:
		pc.endState <- struct{}{}
		return nil
	case 0x30:
		fmt.Println("received publish")
		msg, err := readPublish(pkt)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return pc.messageHandler(msg.Payload)
	}
	return nil
}

func (pc *PorterClient) Publish(topic string, message any) error {
	// TODO implement
	return nil
}
