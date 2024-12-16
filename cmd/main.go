package main

import (
	"net"
	"redis-go/internal/kvstore"
	"redis-go/server"
)

func main() {
	s := kvstore.NewStore()

	listener, err := net.Listen("tcp", ":6379")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go server.HandleConnection(conn, s)
	}
}
