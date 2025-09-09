package persistence

import (
	"encoding/json"
	"fmt"

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
    fmt.Println("[error] Unmarshal error")
		return false
	}
	for _, saved_user := range users {
		if saved_user.Username == user.Username {
    fmt.Println("[error] user already exists")
			return false
		}
	}
	users = append(users, user)
	users_bytes, err := json.Marshal(users)
	if err != nil {
    fmt.Println("[error] Marshall error")
		return false
	}
	_, err = OverwriteFile(filepath, users_bytes)
  fmt.Println("[debug] OverwriteFile error:", err)
	return err == nil
}

func RetrieveUser(filepath string, username string) *share.User {
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return nil
	}
	var users []share.User
	err = json.Unmarshal(f_bytes, &users)
	if err != nil {
		return nil
	}
	for _, saved_user := range users {
		if saved_user.Username == username {
			return &saved_user
		}
	}
	return nil
}

func DeleteUser(filepath string, user share.User) bool {
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return false
	}
	var users []share.User
	err = json.Unmarshal(f_bytes, &users)
	if err != nil {
		return false
	}
	for id, saved_user := range users {
		if saved_user.Username == user.Username && saved_user.Password == user.Password {
      users = append(users[:id], users[id+1:]...) // removing given user only if username and password matches
			users_bytes, err := json.Marshal(users)
			if err != nil {
				return false
			}
			_, err = OverwriteFile(filepath, users_bytes)
      
			return err == nil
		}
	}
	return true // returns true if didn't find user
}

func UpdateUser(filepath string, old_user share.User, new_user share.User) bool {
  ok := DeleteUser(filepath, old_user)
  fmt.Println("[debug] user deleted?", ok)
  if ok {
    ok = SaveUser(filepath, new_user)
    fmt.Println("[debug] user saved?", ok)
    return ok
  }
  return false
}

func RetrieveAllUsers(filepath string) []share.User {
  users := make([]share.User, 0)
	f_bytes, err := ReadFile(filepath)
	if err != nil {
		return users
	}
	json.Unmarshal(f_bytes, &users)
  return  users
}
