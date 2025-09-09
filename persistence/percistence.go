package persistence

import (
	"encoding/json"

	"github.com/vini464/wizard-duel/share"
)

func SaveUser(filepath string, user share.User) bool {
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return false
	}
	var users []share.User
	err = json.Unmarshal(f_bytes, &users)
	if err != nil {
		return false
	}
	for _, saved_user := range users {
		if saved_user.Username == user.Username {
			return false
		}
	}
	users = append(users, user)
	users_bytes, err := json.Marshal(users)
	if err != nil {
		return false
	}
	OverwriteFile(filepath, users_bytes)
	return true
}
