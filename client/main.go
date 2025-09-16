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

var HOSTNAME = "localhost"

func main() {
	for {
		// LOGIN - REGISTER page
		HOSTNAME = Input("Insert server HOSTNAME\n> ")
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
	conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(HOSTNAME, share.SERVERPORT))
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
		fmt.Println("[error] - error connection")
		return false
	}
	return true
}

func login(credentials map[string]string) {
	var user share.User
	conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(HOSTNAME, share.SERVERPORT))
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
	fmt.Println("You are connected")
	conn.Close()
	mainPage(uuid, user)
}

func mainPage(uuid string, user share.User) {
	for {
		choice := Menu("> MAIN PAGE <", "play", "booster", "create deck", "see cards", "see decks", "exit")
		message := share.Message{
			Uuid: uuid,
		}

		switch choice {
		case 0:
			conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(HOSTNAME, share.SERVERPORT))
			if err != nil {
				fmt.Println("[error] - Unable to Connect")
				return
			}
			message.Type = share.PLAY
			err = share.SendMessage(conn, message)
			playing := true
			var self share.ShowableData
			var op share.HiddenData
			for playing {
				share.ReceiveMessage(conn, &message)
				switch message.Type {
				case share.INQUEUE:
					fmt.Println("You are in queue... wait for an oponent")
				case share.WINNER, share.LOSER:
					fmt.Println("Game Over")
					if message.Type == share.WINNER {
						fmt.Println("You win")
					} else {
						fmt.Println("You lose")
					}
					playing = false

				case share.ERROR:
					fmt.Println("an error ocurred!")
					playing = false

				case share.UPDATEGAMESTATE:
					var gamesstate share.GameState
					json.Unmarshal(message.Data, &gamesstate)
					json.Unmarshal(gamesstate.Self, &self)
					json.Unmarshal(gamesstate.Opponent, &op)
					if self.Phase != "WAIT" {
						if self.Phase == "MAIN" {
							hand := make([]string, 0)
							var handcards []share.Card
							json.Unmarshal(self.Hand, &handcards)

							fmt.Println("Game Info")
							fmt.Println("You - energy:", self.Energy, "crystals:", self.Crystals, "life: ", self.HP, "shield:", self.SP, "deck:", self.DeckSize, "DamageBonus: ", self.DamageBonus)
							fmt.Println("Opponent - energy:", op.Energy, "cystals:", op.Crystals, "life: ", op.HP, "shield:", op.SP, "deck:", op.DeckSize, "DamageBonus: ", op.DamageBonus)

							fmt.Println("Your hand:")
							for _, n := range handcards {
								fmt.Println(n)
								if n.Cost <= self.Energy {
									hand = append(hand, n.Name)
								}
							}

							c := Menu("Choose your action:", "place card", "skip phase")

							if len(hand) > 0 && c == 0 {
								cardid := Menu("Choose a card: ", hand...)
								message := share.Message{Type: share.PLACECARD, Data: []byte(hand[cardid])}
								share.SendMessage(conn, message)
							} else {
								message := share.Message{Type: share.SKIPPHASE}
								share.SendMessage(conn, message)
							}

						} else {
							message := share.Message{Type: share.SKIPPHASE}
							share.SendMessage(conn, message)
						}
					}
				case share.SELECTDECK:
					fmt.Println("Select deck:")
					names := make([]string, 0)
					for name := range user.Decks {
						names = append(names, name)
					}
					choice := Menu("Choose a deck", names...)
					message.Type = share.OK
					message.Data = []byte(names[choice])
					share.SendMessage(conn, message)
				case share.OPONENTMOVE:
					fmt.Println("Oponent played:", string(message.Data))
					fmt.Println("Game Info")
					fmt.Println("You - energy:", self.Energy, "crystals:", self.Crystals, "life: ", self.HP, "shield:", self.SP, "deck:", self.DeckSize, "DamageBonus: ", self.DamageBonus)
					fmt.Println("Opponent - energy:", op.Energy, "cystals:", op.Crystals, "life: ", op.HP, "shield:", op.SP, "deck:", op.DeckSize, "DamageBonus: ", op.DamageBonus)
				case share.OPONENTNAME:
					fmt.Println("Your are playing against:", string(message.Data))
				}
			}

			conn.Close()
		case 1:
			conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(HOSTNAME, share.SERVERPORT))
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
				conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(HOSTNAME, share.SERVERPORT))
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
			fmt.Println("Your cards:")
			for _, card := range user.Cards {
				fmt.Print(card.Name, "| cost:", card.Cost, "| rariry:", card.Rarity, "| type:", card.Type, "| effect: ")
				for _, ef := range card.Effects {
					fmt.Print(ef.Type, "-", ef.Amount)
				}
				fmt.Println("")
			}
		case 4:
			fmt.Println("Your decks:")
			for name, deck := range user.Decks {
				fmt.Println(name, ":")
				for _, card := range deck {
					fmt.Print(card.Name, "| cost:", card.Cost, "| rariry:", card.Rarity, "| type:", card.Type, "| effect: ")
					for _, ef := range card.Effects {
						fmt.Print(ef.Type, "-", ef.Amount)
					}
					fmt.Println("")
				}
			}
		case 5:
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
