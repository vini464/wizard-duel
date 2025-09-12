package game

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net"

	"github.com/vini464/wizard-duel/share"
)

// this structure is used only for the game manager function (server-side)
type PlayerGameData struct {
	HP          int
	SP          int
	Energy      int
	Crystals    int
	DamageBonus int
	Username    string
	Phase       string
	Hand        []share.Card
	Deck        []share.Card
	Graveyard   []share.Card
}

const (
	DRAW        = "DRAW"
	REFILL      = "REFILL"
	MAIN        = "MAIN"
	MAINTENANCE = "MAINTENANCE"
	WAIT        = "WAIT"
)

// This function handles with the game logic, returning the name of the winner
func GameManagement(player_conn net.Conn, opponent_conn net.Conn, player_deck []share.Card, opponent_deck []share.Card, player_username string, opponent_username string) string {
	player_data := NewPlayerGameData(player_deck, player_username)
	opponnent_data := NewPlayerGameData(opponent_deck, opponent_username)

	message := share.Message{
		Type: share.OK,
		Data: []byte(player_username),
	}

	share.SendMessage(opponent_conn, message)
	message.Data = []byte(opponent_username)
	share.SendMessage(player_conn, message)

	player_message := make(chan share.Message)
	opponent_message := make(chan share.Message)

	pl_stop := false
	op_stop := false

	go receiving(player_conn, player_message, &pl_stop)
	go receiving(opponent_conn, opponent_message, &op_stop)

	// handling
	for player_data.HP > 0 && opponnent_data.HP > 0 {
		select {

		case player_mes := <-player_message:
			if player_data.Phase != WAIT {
				fmt.Println("player", player_mes)
				switch player_mes.Type {

				case share.ERROR:
					fmt.Println("Some shit Happened", string(player_mes.Data))
					op_stop = true // ends opponent goroutine
					return opponent_username

				case share.SURRENDER:
					winner(opponent_conn)
					loser(player_conn)
					return opponent_username

				case share.PLACECARD:
					cardname := string(message.Data) // receives the name of the card
					for id, card := range player_data.Hand {
						if cardname == card.Name {
							HandleCard(&player_data, &opponnent_data, card)
							player_data.Graveyard = append(player_data.Graveyard, card)
							player_data.Hand = append(player_data.Hand[:id], player_data.Hand[id+1:]...)
							break
						}
					}
				case share.SKIPPHASE:
					HandlePhase(&player_data, &opponnent_data)

				}
			}
		case opponent_mes := <-opponent_message:
			if opponnent_data.Phase != WAIT {

				fmt.Println("Some shit Happened", string(opponent_mes.Data))
				switch opponent_mes.Type {

				case share.ERROR:
					fmt.Println("Some shit Happened")
					pl_stop = true // ends opponent goroutine
					return player_username

				case share.SURRENDER:
					winner(player_conn)
					loser(opponent_conn)
					return player_username

				case share.PLACECARD:
					cardname := string(message.Data) // receives the name of the card
					for id, card := range opponnent_data.Hand {
						if cardname == card.Name {
							HandleCard(&opponnent_data, &player_data, card)
							opponnent_data.Graveyard = append(opponnent_data.Graveyard, card)
							opponnent_data.Hand = append(opponnent_data.Hand[:id], opponnent_data.Hand[id+1:]...)
							break
						}
					}
				case share.SKIPPHASE:
					HandlePhase(&opponnent_data, &player_data)
				}
			}
		}

		// getting game state to send messages
		b_pl_show, _ := json.Marshal(GetShowableData(player_data))
		b_pl_hidd, _ := json.Marshal(GetHiddenData(player_data))
		b_op_show, _ := json.Marshal(GetShowableData(opponnent_data))
		b_op_hidd, _ := json.Marshal(GetHiddenData(opponnent_data))
		op_gamestate := share.GameState{
			Self:     b_op_show,
			Opponent: b_pl_hidd,
		}
		pl_gamestate := share.GameState{
			Self:     b_pl_show,
			Opponent: b_op_hidd,
		}

		message.Type = share.UPDATEGAMESTATE
		// send message to player
		message.Data, _ = json.Marshal(pl_gamestate)
		share.SendMessage(player_conn, message)

		// send message to opponent
		message.Data, _ = json.Marshal(op_gamestate)
		share.SendMessage(opponent_conn, message)
	}

	// Returning who won the game
	if player_data.HP <= 0 {
		winner(opponent_conn)
		loser(player_conn)
		return opponent_username
	}
	winner(player_conn)
	loser(opponent_conn)
	return player_username
}

func winner(conn net.Conn) {
	msg := share.Message{
		Type: share.WINNER,
		Data: []byte("You Win"),
	}
	share.SendMessage(conn, msg)
}
func loser(conn net.Conn) {
	msg := share.Message{
		Type: share.LOSER,
		Data: []byte("You lose"),
	}
	share.SendMessage(conn, msg)
}

func receiving(conn net.Conn, channel chan share.Message, stop *bool) {
	var message share.Message
	for !(*stop) {
		err := share.ReceiveMessage(conn, &message)
		if err != nil {
			message.Type = share.ERROR
			message.Data = []byte(err.Error())
			return // kill function if an error occurred
		}
		channel <- message
	}
}

func NewPlayerGameData(deck []share.Card, username string) PlayerGameData {
	return PlayerGameData{
		HP:          15,
		SP:          10,
		Energy:      0,
		Crystals:    0,
		DamageBonus: 0,
		Username:    username,
		Deck:        shuffle(deck),
		Hand:        make([]share.Card, 0),
		Graveyard:   make([]share.Card, 0),
	}
}

func GetHiddenData(playerData PlayerGameData) share.HiddenData {
	grave, err := json.Marshal(playerData.Graveyard)
	if err != nil {
		grave = nil
	}
	return share.HiddenData{
		HP:          playerData.HP,
		SP:          playerData.SP,
		Energy:      playerData.Energy,
		Crystals:    playerData.Crystals,
		DamageBonus: playerData.DamageBonus,
		Username:    playerData.Username,
		DeckSize:    len(playerData.Deck),
		HandSize:    len(playerData.Hand),
		Graveyard:   grave,
	}
}

func GetShowableData(playerData PlayerGameData) share.ShowableData {
	grave, err := json.Marshal(playerData.Graveyard)
	if err != nil {
		grave = nil
	}
	hand, err := json.Marshal(playerData.Hand)
	if err != nil {
		grave = nil
	}
	return share.ShowableData{
		HP:          playerData.HP,
		SP:          playerData.SP,
		Energy:      playerData.Energy,
		Crystals:    playerData.Crystals,
		DamageBonus: playerData.DamageBonus,
		DeckSize:    len(playerData.Deck),
		Hand:        hand,
		Graveyard:   grave,
		Username:    playerData.Username,
	}
}

func shuffle(deck []share.Card) []share.Card {
	perm := rand.Perm(len(deck))
	for i := range deck {
		j := perm[i]
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func HandleCard(player *PlayerGameData, opponent *PlayerGameData, card share.Card) {
	for _, effect := range card.Effects {
		switch effect.Type {
		case "heal":
			player.HP += effect.Amount
		case "shield":
			player.SP += effect.Amount
		case "damage":
			damage := effect.Amount + player.DamageBonus
			player.DamageBonus = 0
			if opponent.SP == 0 {
				opponent.HP -= damage
			} else if opponent.SP < damage {
				opponent.HP -= damage - opponent.SP
				opponent.SP = 0
			} else {
				opponent.SP -= damage
			}
		case "energy":
			player.Crystals += effect.Amount
		case "destroy_enemy_shield":
			opponent.SP = 0
		case "next_spell_damage_bonus":
			player.DamageBonus += effect.Amount
		case "draw":
			if effect.Amount <= len(player.Deck) {
				player.Hand = append(player.Hand, player.Deck[:effect.Amount]...)
			} else {
				player.Hand = append(player.Hand, player.Deck...)
			}
		case "discard":
			opponent.Hand = shuffle(opponent.Hand)
			opponent.Graveyard = append(opponent.Graveyard, opponent.Hand[:effect.Amount-1]...)
			opponent.Hand = opponent.Hand[effect.Amount:]
		}
	}
}

func HandlePhase(player *PlayerGameData, opponent *PlayerGameData) {
	if player.Phase == DRAW {
		player.Hand = append(player.Hand, player.Deck[0])
		player.Deck = player.Deck[1:]
		player.Phase = REFILL
	} else {
		actualPhase := skipPhase(player.Phase)
		switch actualPhase {
		case REFILL:
			player.Crystals++
			player.Energy += player.Crystals
		case MAINTENANCE:
			if len(player.Hand) > 6 {
				player.Hand = shuffle(player.Hand)
				player.Hand = player.Hand[(len(player.Hand) - 6):]
			}
		case WAIT:
			opponent.Phase = DRAW
		}
	}
}

func skipPhase(actualPhase string) string {
	switch actualPhase {
	case DRAW:
		return REFILL
	case REFILL:
		return MAIN
	case MAIN:
		return MAINTENANCE
	case MAINTENANCE:
		return WAIT
	case WAIT:
		return DRAW
	default:
		return MAIN
	}
}
