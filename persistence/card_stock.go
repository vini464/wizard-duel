package persistence

import (
	"encoding/json"
	"math/rand/v2"

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

func RemoveFromStock(cards ...share.Card) {
	stock := RetrieveStock()
	for _, card := range cards {
		for id, card_stock := range stock {
			if card_stock.Card.Name == card.Name {
				stock[id].Quantity--
			}
		}
	}
}

func AddToStock(cards ...share.Card) {
	stock := RetrieveStock()
	for _, card := range cards {
		for id, card_stock := range stock {
			if card_stock.Card.Name == card.Name {
				stock[id].Quantity++
			}
		}
	}
}

func CreateBooster() []share.Card {
	booster := make([]share.Card, 0)

	common_cards := GetByRarity("common")
	uncommon_cards := GetByRarity("uncommon")
	rare_cards := GetByRarity("rare")
	legendary_cards := GetByRarity("legendary")

	for range 3 { // common cards
		r := rand.IntN(len(common_cards))
		booster = append(booster, common_cards[r])
		common_cards = append(common_cards[:r], common_cards[r+1:]...)
	}

	// uncommon cards
	booster = append(booster, common_cards[rand.IntN(len(uncommon_cards))])

	// rare or legendary cards
	r := rand.IntN(100)
	if r < 15 && len(legendary_cards) > 0 {
		booster = append(booster, legendary_cards[rand.IntN(len(legendary_cards))])
	} else {
		booster = append(booster, rare_cards[rand.IntN(len(rare_cards))])
	}

  RemoveFromStock(booster...) // remove all used cards from stock

	return booster
}



func InitializeStock() {
	cards := RetrieveAllCards(CARDTYPES)
	for _, card := range cards {
		AddCardToStock(card)
	}
}
