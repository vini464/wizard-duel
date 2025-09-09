package game

import (
	"sync"

	"github.com/vini464/wizard-duel/share"
)


type PlayerGameData struct {
	Username    string
	Hand        []share.Card
	Deck        []share.Card
	Graveyard   []share.Card
	HP          int
	SP          int
	Energy      int
	Crystals    int
	DamageBonus int
}

type PrivateGameState struct {
	Mutex       *sync.Mutex
	PlayersData map[string]PlayerGameData // map entre o username e os dados
	Turn        string
	Phase       string
	Round       int
}

