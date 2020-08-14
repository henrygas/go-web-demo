package tcpdemo

import (
	"fmt"
	"log"
	"net"
)

type TcpServer struct {
	network string
	address string
}

func NewTcpServer() (*TcpServer, error) {
	return &TcpServer{
		network: "tcp",
		address: address,
	}, nil
}

func (ts *TcpServer) Start() error {
	fmt.Println("[Server]Starting tcp server...")
	if listener, err := net.Listen(ts.network, ts.address); err != nil {
		log.Println("[Server]failed to start listen tcp, err: ", err)
		return err
	} else {
		fmt.Printf("[Server]succ to start tcp server, listening on %s %s\n", ts.network, ts.address)
		for {
			if conn, err := listener.Accept(); err != nil {
				log.Println("[Server]failed to accept listener, err: ", err)
				return err
			} else {
				go ts.serve(conn)
			}
		}
	}
}

func (ts *TcpServer) serve(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		if n, err := conn.Read(buf); err != nil {
			log.Println("[Server]failed to read from conn to buf, err: ", err)
			return
		} else {
			fmt.Println("[Server]succ to receive data: ", string(buf[:n]))
		}
	}
}