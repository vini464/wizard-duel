package share

import (
	"crypto/md5"
	"encoding/hex"
)

func HashText(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum([]byte("secure text")))
}
