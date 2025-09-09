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
	conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
	if err != nil {
    fmt.Println("[error] - Unable to Connect")
		return
	}
	defer conn.Close()
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

  fmt.Println("[debug] - Connected, uuid:", uuid)
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
