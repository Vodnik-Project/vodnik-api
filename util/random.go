package util

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomNumber(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	var randomstring strings.Builder
	k := len(letters)
	for i := 0; i < n; i++ {
		randomstring.WriteByte(letters[rand.Intn(k)])
	}
	return randomstring.String()
}

func RandomEmail() string {
	return RandomString(5) + "@" + RandomString(4) + ".com"
}
