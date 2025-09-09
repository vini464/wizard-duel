package persistence

import (
	"encoding/json"

	"github.com/vini464/wizard-duel/share"
)


func SaveCard(filepath string, card share.Card) bool {
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return false
	}
	var cards []share.Card
	err = json.Unmarshal(f_bytes, &cards)
	if err != nil {
		return false
	}
	for _, saved_card := range cards {
		if saved_card.Name == card.Name {
			return false
		}
	}
	cards = append(cards, card)
	users_bytes, err := json.MarshalIndent(cards, "", " ")
	if err != nil {
		return false
	}
	_, err = OverwriteFile(filepath, users_bytes)
	return err == nil
}

func RetrieveCard(filepath string, name string) *share.Card {
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return nil
	}
  var cards []share.Card
	err = json.Unmarshal(f_bytes, &cards)
	if err != nil {
		return nil
	}
	for _, saved_card := range cards {
		if saved_card.Name == name {
			return &saved_card
		}
	}
	return nil
}

func DeleteCard(filepath string, card share.Card) bool {
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return false
	}
	var cards []share.Card
	err = json.Unmarshal(f_bytes, &cards)
	if err != nil {
		return false
	}
	for id, saved_card := range cards {
		if saved_card.Name == card.Name {
      cards = append(cards[:id], cards[id+1:]...) // removing given card only if name matches
			cards_bytes, err := json.MarshalIndent(cards, "", " ")
			if err != nil {
				return false
			}
			_, err = OverwriteFile(filepath, cards_bytes)
			return err == nil
		}
	}
	return true // returns true if didn't find card
}

func UpdateCard(filepath string, old_card share.Card, new_card share.Card) bool {
  ok := DeleteCard(filepath, old_card)
  if ok {
    ok = SaveCard(filepath, new_card)
    return ok
  }
  return false
}

func RetrieveAllCards(filepath string) []share.Card {
  cards := make([]share.Card, 0)
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return cards
	}
	json.Unmarshal(f_bytes, &cards)
  return cards
}
