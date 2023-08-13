package util

import (
	"math/rand"
	"strings"
	"time"
)

const alph = "abcdefghijklmnoprstvxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RadnomString(n int) string {
	var sb strings.Builder
	k := len(alph)

	for i := 0; i < n; i++ {
		c := alph[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomOwner() string {
	return RadnomString(6)
}

func RadnomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD, HRK}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
