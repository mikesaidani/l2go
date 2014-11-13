package gameserver

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {

	fmt.Println("A client is trying to connect to the Game Server...")
	conn.Close()

}

func Init() {
	ln, err := net.Listen("tcp", ":7777")
	defer ln.Close()

	if err != nil {
		fmt.Println("Couldn't initialize the Game Server")
	} else {
		fmt.Println("Game Server initialized.")
		fmt.Println("Listening on 127.0.0.1:7777.")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Couldn't accept the incoming connection.")
			continue
		} else {
			go handleConnection(conn)
		}

	}
}
