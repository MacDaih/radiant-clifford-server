package tcpclient

import (
    "log"
	"net"
	"time"
)

type readerFunc func(net.Conn) error

func RunTCPCLient(socket string, key string, r readerFunc) error {
	log.Println("Running Collector target")
	addr, err := net.ResolveTCPAddr("tcp4", socket)
	if err != nil {
		log.Printf("resolving address error : %s\n", err.Error())
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		log.Printf("dial to address error : %s\n", err.Error())
		return err
	}

	defer conn.Close()

	for {
		_, err := conn.Write([]byte(key))
		if err != nil {
			log.Println(err)
			continue
		}
		r(conn)
		time.Sleep(time.Second * 10)
	}
}
