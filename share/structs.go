package share

// this package only have structs type for server and client usage

type User struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Coins    int               `json:"coins"`
	Cards    []Card            `json:"cards"`
	Decks    map[string][]Card `json:"decks"` // deck name -> card list
}

type Card struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Rarity  string   `json:"rarity"`
	Cost    int      `json:"cost"`
	Effects []Effect `json:"effects"`
}

type Effect struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

func NewUser(username string, password string) *User {
	if username == "" || password == "" {
		return nil
	}
	return &User{
		Username: username,
		Password: HashText(password),
		Coins:    120,
		Cards:    make([]Card, 0),
		Decks:    make(map[string][]Card),
	}
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
	Phase       string `json:"phase"`
}

// Information that the player have access
type ShowableData struct {
	HP          int    `json:"hp"`
	SP          int    `json:"sp"`
	Energy      int    `json:"energy"`
	Crystals    int    `json:"crystals"`
	DamageBonus int    `json:"damage-bonus"`
	DeckSize    int    `json:"deck-size"`
	Hand        []byte `json:"hand"`
	Graveyard   []byte `json:"graveyard"`
	Username    string `json:"username"`
	Phase       string `json:"phase"`
}

type GameState struct {
	Self     []byte `json:"self"`
	Opponent []byte `json:"opponent"`
}
