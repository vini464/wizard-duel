package main

import (
	"fmt"
	"io"
	"net"

	"github.com/vini464/wizard-duel/share"
)

// replace this for persistence later
var USERS = make([]share.User, 0)

func main() {
	server, err := net.Listen("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("[error] - connection lost")
		} else {
			go handle_client(conn)
		}
	}
}

func handle_client(conn net.Conn) {
	var message share.Message
	err := share.ReceiveMessage(conn, &message)
	if err != nil {
		if err == io.EOF {
			fmt.Println("[info] - Client disconnected!")
		} else {
			fmt.Println("[error] - Connection lost!")
		}
	}
	switch message.Type {
	case share.REGISTER:
		register(message.Data)
	case share.LOGIN:
		login(message.Data)
	case share.SAVEDECK:
		save_deck(message.Data)
	case share.GETBOOSTER:
		get_booster(message.Data)
	case share.PLAY:
		play(message.Data)
	default:
		fmt.Println("[error] unknow command:", message.Type)
	}
}

func register(serialized_data []byte) {
  

}
func login(serialized_data []byte) {

}
func save_deck(serialized_data []byte) {

}
func get_booster(serialized_data []byte) {

}
func play(serialized_data []byte) {

}
