package common

import (
	"dexbg/types"
	"time"
)

func GenerateKey() types.Key {
	var key types.Key
	key.Key = RandomString(15)
	now := time.Now()
	key.Expiration = int(now.Unix()) + 345600
	return key
}
