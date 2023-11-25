package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" //áéíóúüñÑÁÉÍÓÚÜ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 9999)
}
func RandomCurrency() string {
	currencies := getCurrencies()
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
func RandomEmail() string {
	return RandomString(int(RandomInt(4, 10))) + "@email.com"
}
