package main

import (
	"go-web/httpdemo"
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

func StartNewHttpServer() {
	if hs, err := httpdemo.NewHttpServer(); err != nil {
		log.Println("failed to NewHttpServer(), err: ", err)
		return
	} else if err := hs.Start(); err != nil {
		log.Println("failed to hs.Start(), err: ", err)
		return
	}
}

func main() {
	StartNewHttpServer()
	//go StartTcpServer()
	//StartTcpClient()
}
