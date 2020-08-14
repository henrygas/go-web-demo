package main

import (
	"go-web/tcpdemo"
	"log"
)

func StartTcpServer() {
	if ts, err := tcpdemo.NewTcpServer(); err != nil {
		log.Println("failed to NewTcpServer(), err: ", err)
		return
	} else {
		if err := ts.Start(); err != nil {
			log.Println("failed to ts.Start(), err: ", err)
			return
		}
	}
}

func StartTcpClient() {
	if tc, err := tcpdemo.NewTcpClient(); err != nil {
		log.Println("failed to NewTcpClient(), err: ", err)
		return
	} else {
		if err := tc.Start(); err != nil {
			log.Println("failed to tc.Start(), err: ", err)
			return
		}
	}
}

func main() {
	go StartTcpServer()
	StartTcpClient()
}