package share

// this package only have structs type for server and client usage

type User struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Cards    []Card            `json:"cards"`  
	Decks    map[string][]Card `json:"decks"` // deck name -> card list
	Coins    int               `json:"coins"`
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
	return &User{
		Username: username,
		Password: password, // hash it later
		Coins:    120,
	}
}
