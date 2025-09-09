package main

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"
	"sync"

	"github.com/vini464/wizard-duel/persistence"
	"github.com/vini464/wizard-duel/share"
)

const (
	USERSFILEPATH = "database/users.json"
	CARDSFILEPATH = "database/cards.json"
)

var ONLINEUSERS = make(map[string]string) // string uuid - string username

func main() {
	server, err := net.Listen(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
	if err != nil {
		panic(err)
	}
	defer server.Close()

	var user_db sync.Mutex
	var card_db sync.Mutex

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("[error] - connection lost")
		} else {
			go handle_client(conn, &user_db, &card_db)
		}
	}
}

func handle_client(conn net.Conn, user_db *sync.Mutex, card_db *sync.Mutex) {
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
		register(message, conn, user_db)
	case share.LOGIN:
		login(message, conn, user_db)
	case share.SAVEDECK:
		save_deck(message, conn, user_db)
	case share.GETBOOSTER:
		get_booster(message, conn, user_db)
	case share.PLAY:
		play(message)
	case share.LOGOUT:
		log_out(message, conn)
	default:
		fmt.Println("[error] unknow command:", message.Type)
	}
}

func register(message share.Message, conn net.Conn, user_db *sync.Mutex) {
	serialized_data := message.Data
	var credentials = make(map[string]string)
	err := json.Unmarshal(serialized_data, &credentials)
	response := share.Message{}
	if err != nil {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	newuser := share.NewUser(credentials["username"], credentials["password"])
	user_db.Lock()
	success := persistence.SaveUser(USERSFILEPATH, *newuser)
	user_db.Unlock()
	if success {
		response.Type = share.OK
		share.SendMessage(conn, response)
		return
	}
	response.Type = share.ERROR
	share.SendMessage(conn, response)

}

func login(message share.Message, conn net.Conn, user_db *sync.Mutex) {
	serialized_data := message.Data
	var credentials = make(map[string]string)
	response := share.Message{}
	err := json.Unmarshal(serialized_data, &credentials)
	if err != nil {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user_db.Lock()
	saved_user := persistence.RetrieveUser(USERSFILEPATH, credentials["username"])
	user_db.Unlock()
	if saved_user != nil {
		if saved_user.Password == credentials["password"] {
			uuid := share.HashText(saved_user.Username)
			ONLINEUSERS[uuid] = saved_user.Username
			response.Type = share.OK
			response.Uuid = uuid
			share.SendMessage(conn, response)
			return
		}
	}
	response.Type = share.ERROR
	share.SendMessage(conn, response)
}

func save_deck(message share.Message, conn net.Conn, user_db *sync.Mutex) {
	serialized_data := message.Data
	var deck = make(map[string][]share.Card) // string deckname - []card cards
	err := json.Unmarshal(serialized_data, &deck)
	response := share.Message{}
	if err != nil {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	username, ok := ONLINEUSERS[message.Uuid]
	if !ok {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user_db.Lock()
	user := persistence.RetrieveUser(USERSFILEPATH, username)
	maps.Copy(user.Decks, deck)
	success := persistence.UpdateUser(USERSFILEPATH, *user, *user)
	user_db.Unlock()
	if success {
		response.Type = share.OK
		share.SendMessage(conn, response)
		return
	}
	response.Type = share.ERROR
	share.SendMessage(conn, response)
}

func get_booster(message share.Message, conn net.Conn, user_db *sync.Mutex) {
	username, ok := ONLINEUSERS[message.Uuid]
	response := share.Message{}
	if !ok {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user_db.Lock()
	user := persistence.RetrieveUser(USERSFILEPATH, username)
	if user.Coins >= 5 { // TODO: create booster logic
		response.Type = share.OK
		share.SendMessage(conn, response)
	}
	response.Type = share.ERROR
	share.SendMessage(conn, response)
}

func log_out(message share.Message, conn net.Conn) {
	_, ok := ONLINEUSERS[message.Uuid]
	response := share.Message{}
	if ok {
		delete(ONLINEUSERS, message.Uuid)
	}
	response.Type = share.OK
	share.SendMessage(conn, response)
}

func play(message share.Message) bool {
	fmt.Println("[debug] - play command", message)
	return true
}
