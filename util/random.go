package util

import (
	"math/rand"
)

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
