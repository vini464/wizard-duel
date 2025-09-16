package main

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"
	"sync"

	"github.com/vini464/wizard-duel/persistence"
	"github.com/vini464/wizard-duel/game"
	"github.com/vini464/wizard-duel/share"
)

const (
	USERSFILEPATH = "database/users.json"
)

type Userinfo struct {
	conn     net.Conn
	username string
}

var ONLINEUSERS = make(map[string]string) // string uuid - string username
var GAMEQUEUE = make([]Userinfo, 0)       // queue that holds user connection

func main() {
	persistence.InitializeStock()
	server, err := net.Listen(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
	if err != nil {
		panic(err)
	}
	defer server.Close()

	var user_db sync.Mutex
	var card_db sync.Mutex
	var queue_mutex sync.Mutex

	fmt.Println("[debug] - Server Ready")
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("[error] - connection lost")
		} else {
			fmt.Println("[debug] - Client Connected", conn.RemoteAddr())
			go handle_client(conn, &user_db, &card_db, &queue_mutex)
		}
	}
}

func handle_client(conn net.Conn, user_db *sync.Mutex, card_db *sync.Mutex, queue_mutex *sync.Mutex) {
	var message share.Message
	err := share.ReceiveMessage(conn, &message)
	fmt.Println("[debug] - error:", err)
	if err != nil {
		if err == io.EOF {
			fmt.Println("[info] - Client disconnected!")
			return
		} else {
			fmt.Println("[error] - Connection lost!")
			return
		}
	}
	switch message.Type {
	case share.REGISTER:
		fmt.Println("[debug] REGISTER command:", message.Type)
		register(message, conn, user_db)
	case share.LOGIN:
		fmt.Println("[debug] LOGIN command:", message.Type)
		login(message, conn, user_db)
	case share.SAVEDECK:
		fmt.Println("[debug] SAVEDECK command:", message.Type)
		save_deck(message, conn, user_db)
	case share.GETBOOSTER:
		fmt.Println("[debug] GETBOOSTER command:", message.Type)
		get_booster(message, conn, user_db, card_db)
	case share.PLAY:
		fmt.Println("[debug] PLAY command:", message.Type)
		play(conn, queue_mutex, ONLINEUSERS[message.Uuid], user_db)
	case share.LOGOUT:
		fmt.Println("[debug] LOGOUT command:", message.Type)
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
		fmt.Println("[error] - error while deserializing")
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
		if saved_user.Password == share.HashText(credentials["password"]) {
			ser, err := json.Marshal(*saved_user)
			if err == nil {
				uuid := share.HashText(saved_user.Username)
				ONLINEUSERS[uuid] = saved_user.Username
				response.Type = share.OK
				response.Uuid = uuid
				response.Data = ser
				share.SendMessage(conn, response)
				return
			}
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

func get_booster(message share.Message, conn net.Conn, user_db *sync.Mutex, card_db *sync.Mutex) {
	username, ok := ONLINEUSERS[message.Uuid]
	response := share.Message{}
	if !ok {
		response.Type = share.ERROR
		share.SendMessage(conn, response)
		return
	}
	user_db.Lock()
	user := persistence.RetrieveUser(USERSFILEPATH, username)
	user_db.Unlock()
	if user.Coins >= 5 {
		user.Coins -= 5
		card_db.Lock()
		booster := persistence.CreateBooster()
		card_db.Unlock()
		ser, err := json.Marshal(booster)
		if err == nil {
			user_db.Lock()
			user.Cards = append(user.Cards, booster...)
			ok := persistence.UpdateUser(USERSFILEPATH, *user, *user)
			user_db.Unlock()
			if ok {
				response.Type = share.OK
				response.Data = ser
				share.SendMessage(conn, response)
				return
			}
		}
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

func play(conn net.Conn, queue_mutex *sync.Mutex, username string, user_db *sync.Mutex) {
	info := Userinfo{conn, username}
	queue_mutex.Lock()
	if len(GAMEQUEUE) == 0 {
		GAMEQUEUE = append(GAMEQUEUE, info)
		queue_mutex.Unlock()
		share.SendMessage(conn, share.Message{Type: share.INQUEUE})
	} else {
		op_info := GAMEQUEUE[0]
		GAMEQUEUE = GAMEQUEUE[1:]
		queue_mutex.Unlock()

		pl_deck := make([]share.Card, 0)
		op_deck := make([]share.Card, 0)

		pl_ok := make(chan bool)
		op_ok := make(chan bool)
		go getPlayerDeck(info, &pl_deck, pl_ok, user_db)
		go getPlayerDeck(op_info, &op_deck, op_ok, user_db)
		ok1 := <-pl_ok
		ok2 := <-op_ok

		fmt.Println("[debug] - ok1:", ok1, " ok2:", ok2)
		if ok1 && ok2 {
			winner := game.GameManagement(conn, op_info.conn, pl_deck, op_deck, info.username, op_info.username)
			user := persistence.RetrieveUser(USERSFILEPATH, winner)
			user.Coins += 5
			persistence.UpdateUser(USERSFILEPATH, *user, *user)

		} else {
			msg := share.Message{
				Type: share.ERROR,
				Data: []byte("An error occurred"),
			}
			share.SendMessage(conn, msg)
			share.SendMessage(op_info.conn, msg)
		}
	}
}

func getPlayerDeck(userinfo Userinfo, deck *[]share.Card, ok chan bool, user_db *sync.Mutex) {
	msg := share.Message{
		Type: share.SELECTDECK,
	}
	err := share.SendMessage(userinfo.conn, msg)
	if err == nil {
		err = share.ReceiveMessage(userinfo.conn, &msg)
		if err == nil {
			deckname := string(msg.Data)
			user_db.Lock()
			user := persistence.RetrieveUser(USERSFILEPATH, userinfo.username)
			user_db.Unlock()
			if user == nil {
				ok <- false
				return
			}
			for n, d := range user.Decks {
				if n == deckname {
					*deck = d
					ok <- true
					return
				}
			}
		}
	}
	ok <- false
}
