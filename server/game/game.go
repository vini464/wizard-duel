package game

import (
	"github.com/vini464/wizard-duel/share"
)

// this struct is used only for the game manager function (server-side)
type PlayerGameData struct {
	HP          int
	SP          int
	Energy      int
	Crystals    int
	DamageBonus int
	Username    string
	Hand        []share.Card
	Deck        []share.Card
	Graveyard   []share.Card
}

// Limited information about the enemy
type HiddenData struct {
	HandSize    int    `json:"hand-size"`
	DeckSize    int    `json:"deck-size"`
	Energy      int    `json:"energy"`
	HP          int    `json:"hp"`
	SP          int    `json:"sp"`
	Crystals    int    `json:"crystals"`
	DamageBonus int    `json:"damage-bonus"`
	Username    string `json:"username"`
	Graveyard   []byte `json:"graveyard"` // this is a []Card serialized
}

// Information that the player have access
type SelfData struct {
	HP          int    `json:"hp"`
	SP          int    `json:"sp"`
	Energy      int    `json:"energy"`
	Crystals    int    `json:"crystals"`
	DamageBonus int    `json:"damage-bonus"`
	DeckSize    int    `json:"deck-size"`
	Hand        []byte `json:"hand"`
	Graveyard   []byte `json:"graveyard"`
}
