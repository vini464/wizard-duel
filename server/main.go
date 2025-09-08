package main

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"

	"github.com/vini464/wizard-duel/share"
)

// TODO: replace this for persistence later
var USERS = make([]*share.User, 0)

var ONLINEUSERS = make(map[string]int) // string uuid - int users index

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
		register(message, conn)
	case share.LOGIN:
		login(message, conn)
	case share.SAVEDECK:
		save_deck(message, conn)
	case share.GETBOOSTER:
		get_booster(message, conn)
	case share.PLAY:
		play(message)
	default:
		fmt.Println("[error] unknow command:", message.Type)
	}
}

func register(message share.Message, conn net.Conn) {
	serialized_data := message.Data
	var credentials = make(map[string]string)
	err := json.Unmarshal(serialized_data, &credentials)
	response := share.Message{}
	if err != nil {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	for _, user := range USERS {
		if user.Username == credentials["username"] {
			response.Type = share.ERROR
			share.SendMessage(conn, response)
			return
		}
	}
	newuser := share.NewUser(credentials["username"], credentials["password"])
	USERS = append(USERS, newuser)
	response.Type = share.OK
	share.SendMessage(conn, response)
}

func login(message share.Message, conn net.Conn) {
	serialized_data := message.Data
	var credentials = make(map[string]string)
	response := share.Message{}
	err := json.Unmarshal(serialized_data, &credentials)
	if err != nil {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	for id, user := range USERS {
		if user.Username == credentials["username"] && user.Password == credentials["password"] {
			ONLINEUSERS[user.Username] = id
			response.Type = share.OK
			response.Uuid = user.Username // TODO: replace this later
			share.SendMessage(conn, response)
			return
		}
	}
	response.Type = share.ERROR
	share.SendMessage(conn, response)
}

func save_deck(message share.Message, conn net.Conn) {
	serialized_data := message.Data
	var deck = make(map[string][]share.Card) // string deckname - []card cards
	err := json.Unmarshal(serialized_data, &deck)
	response := share.Message{}
	if err != nil {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user_id, ok := ONLINEUSERS[message.Uuid]
	if !ok {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user := USERS[user_id]
	maps.Copy(user.Decks, deck)
	response.Type = share.OK
	share.SendMessage(conn, response)
}

func get_booster(message share.Message, conn net.Conn) {
	user_id, ok := ONLINEUSERS[message.Uuid]
	response := share.Message{}
	if !ok {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user := USERS[user_id]
	if user.Coins >= 5 {
		response.Type = share.OK
		share.SendMessage(conn, response)

	}
	response.Type = share.ERROR
	share.SendMessage(conn, response)
}

func play(message share.Message) bool {
	fmt.Println("[debug] - play command", message)
	return true
}
