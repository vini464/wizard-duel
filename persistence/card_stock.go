package persistence

import (
	"encoding/json"

	"github.com/vini464/wizard-duel/share"
)

type StockCard struct {
	Card     share.Card `json:"card"`
	Quantity int        `json:"quantity"`
	Min      int        `json:"min"`
}

const (
	CARDSTOCK = "database/cardstock.json" // here i save stockCard structure
	CARDTYPES = "database/cardtypes.json" // here just have one copy of any card
)

func RetrieveStock() []StockCard {
	stock := make([]StockCard, 0)
	f_bytes, err := ReadFile(CARDSTOCK)
	if err == nil {
		json.Unmarshal(f_bytes, &stock)
	}
	return stock
}

func ReplaceStock(new_stock []StockCard) bool {
	ser, err := json.Marshal(new_stock)
	if err == nil {
		_, err = OverwriteFile(CARDSTOCK, ser)
	}
	return err == nil
}

func UpdateStock() {
	stock := RetrieveStock()
	for id, card := range stock {
		if card.Quantity < card.Min {
			stock[id].Quantity = card.Min
		}
	}
	ReplaceStock(stock)
}

func AddCardToStock(card share.Card) {
	var card_min int
	switch card.Rarity {
	case "common":
		card_min = 32
	case "uncommon":
		card_min = 16
	case "rare":
		card_min = 8
	case "legendary":
		card_min = 4
	default:
		card_min = 0
	}
	card_stock := StockCard{
		Card:     card,
		Min:      card_min,
		Quantity: card_min,
	}
	stock := RetrieveStock()
	found := false
	for id, s_card := range stock {
		if s_card.Card.Name == card.Name {
			stock[id].Min = card_min
			found = true
			break
		}
	}
	if !found {
		stock = append(stock, card_stock)
	}
	ReplaceStock(stock)
}

func GetByRarity(rarity string) []share.Card {
	commons := make([]share.Card, 0)
	stock := RetrieveStock()
	for _, card_stock := range stock {
		if card_stock.Card.Rarity == rarity {
			for range card_stock.Quantity {
				commons = append(commons, card_stock.Card)
			}
		}
	}
	return commons
}
