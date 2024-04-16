package string_util

import (
	"golang.org/x/exp/rand"
	"time"
)

var CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(n int) string {
	rand.Seed(uint64(time.Now().UnixNano()))

	var result string
	charsetLen := len(CHARSET)
	for i := 0; i < n; i++ {
		result += string(CHARSET[rand.Intn(charsetLen)])
	}
	return result
}
