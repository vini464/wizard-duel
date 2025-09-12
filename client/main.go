package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/vini464/wizard-duel/share"
)

func main() {
	for {
		// LOGIN - REGISTER page
		credentials := GetCredentials()
		choice := Menu("> Wizard Duel <", "login", "register", "exit")
		switch choice {
		case 0:
			login(credentials)
		case 1:
			ok := register(credentials)
			if ok {
				fmt.Println("[debug] - User registered!")
			} else {
				fmt.Println("[debug] - User not registered!")
			}
		case 2:
			return
		}
	}
}

func register(credentials map[string]string) bool {
	conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
	if err != nil {
		fmt.Println("[error] - Unable to Connect")
		return false
	}
	defer conn.Close()
	ser, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("[error] - Serialization Failed")
		return false
	}
	message := share.Message{
		Type: share.REGISTER,
		Data: ser,
	}
	err = share.SendMessage(conn, message)
	fmt.Println("[error] - error connection")
	if err != nil {
		return false
	}
	err = share.ReceiveMessage(conn, &message)
	if err != nil || message.Type == share.ERROR {
		fmt.Println("[error] - error connection, message type:", message.Type)
		return false
	}
	return true
}

func login(credentials map[string]string) {
	var user share.User
	conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
	if err != nil {
		fmt.Println("[error] - Unable to Connect")
		return
	}
	ser, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("[error] - Unable to Connect")
		return
	}
	message := share.Message{
		Type: share.LOGIN,
		Data: ser,
	}
	err = share.SendMessage(conn, message)
	if err != nil {
		fmt.Println("[error] - Unable to Connect")
		return
	}
	err = share.ReceiveMessage(conn, &message)
	if err != nil {
		fmt.Println("[error] - Unable to Connect")
		return
	}
	if message.Type == share.ERROR {
		fmt.Println("[error] - wrong username or password")
		return
	}
	uuid := message.Uuid
	err = json.Unmarshal(message.Data, &user)
	if err != nil {
		fmt.Println("[error] - Unable to Connect")
		return
	}

	fmt.Println("[debug] - Connected, uuid:", uuid)
	fmt.Println("[debug] - Connected, user:", user)
	conn.Close()
	mainPage(uuid, user)
}

func mainPage(uuid string, user share.User) {
	for {
		choice := Menu("> MAIN PAGE <", "play", "booster", "create deck", "exit")
		message := share.Message{
			Uuid: uuid,
		}

		switch choice {
		case 0:
			conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
			if err != nil {
				fmt.Println("[error] - Unable to Connect")
				return
			}
			message.Type = share.PLAY
      err = share.SendMessage(conn, message)
      err = share.ReceiveMessage(conn, &message)
      fmt.Println("[debug] - You are playing with:", string(message.Data))
      conn.Close()
		case 1:
			conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
			if err != nil {
				fmt.Println("[error] - Unable to Connect")
				return
			}
			fmt.Println("[debug] - getting booster")
			message.Type = share.GETBOOSTER
			err = share.SendMessage(conn, message)
			if err != nil {
				fmt.Println("[error] - Unable to Connect")
				return
			}
			err = share.ReceiveMessage(conn, &message)
			if err != nil {
				fmt.Println("[error] - Unable to Connect")
				return
			}
			if message.Type == share.OK {
				var booster []share.Card
				err = json.Unmarshal(message.Data, &booster)
				if err != nil {
					fmt.Println("[error] - Unmarshal error")
					return
				}
				println("[debug] - booster cards:", booster)
				user.Cards = append(user.Cards, booster...)
			} else {
				println("[debug] - Response", message.Type)
			}
    conn.Close()
		case 2:
			if len(user.Cards) > 20 {
				deckcards := make([]share.Card, 0)
				deckname := Input("Deck name: ")
				for len(deckname) < 5 {
					deckname = Input("Deck name: ")
				}
				for len(deckcards) < 20 {
					for id, card := range user.Cards {
						fmt.Println(id, "-", card)
					}
					choice, err := strconv.Atoi(Input("choose a card: "))
					if err == nil && choice >= 0 && choice < len(user.Cards) {
						deckcards = append(deckcards, user.Cards[choice])
						user.Cards = append(user.Cards[:choice], user.Cards[choice+1:]...)
					}
				}
				deck := make(map[string][]share.Card)
				deck[deckname] = deckcards
				ser, err := json.Marshal(deck)
				if err != nil {
					return
				}
				conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
				if err != nil {
					fmt.Println("[error] - Unable to Connect")
					return
				}
				message.Data = ser
				message.Type = share.SAVEDECK
				err = share.SendMessage(conn, message)
				if err != nil {
					return
				}
				err = share.ReceiveMessage(conn, &message)
				fmt.Println("[debug] response was:", message.Type)
        conn.Close()
			}
		case 3:
			return
		}
	}
}

func Menu(title string, args ...string) int {
	for {
		fmt.Println(title)
		for id, op := range args {
			fmt.Println(id+1, "-", op+";")
		}
		input := Input("Select an option:\n> ")
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("[error] - Invalid Option!!")
		}
		if choice > 0 && choice <= len(args) {
			return choice - 1
		}
	}
}

func Input(title string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(title)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func GetCredentials() map[string]string {
	username := Input("username: ")
	for len(username) < 3 {
		fmt.Println("[error] - invalid username")
		username = Input("username: ")
	}
	password := Input("password: ")
	for len(password) < 3 {
		fmt.Println("[error] - invalid password")
		password = Input("password: ")
	}
	data := make(map[string]string)
	data["username"] = username
	data["password"] = password
	return data
}
