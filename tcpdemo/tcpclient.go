package tcpdemo

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type TcpClient struct {
	conn net.Conn
}

func NewTcpClient() (*TcpClient, error) {
	tc := &TcpClient{}
	if conn, err := net.Dial("tcp", address); err != nil {
		log.Println("failed to Dial tcp server, err: ", err)
		return nil, err
	} else {
		tc.conn = conn
	}

	return tc, nil
}

func (tc *TcpClient) Start() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("[Client]What's your name?")
	clientName, _ := reader.ReadString('\n')
	trimmedClientName := strings.Trim(clientName, "\r\n")

	for {
		fmt.Println("[Client]What do you want to send to server?")
		content, _ := reader.ReadString('\n')
		trimmedContent := strings.Trim(content, "\r\n")

		if _, err := tc.conn.Write([]byte(content)); err != nil {
			log.Println("[Client]failed to write content to server, err: ", err)
			return err
		} else {
			fmt.Printf("[Client]%s says %s to tcp server\n", trimmedClientName, trimmedContent)
		}
	}
}
