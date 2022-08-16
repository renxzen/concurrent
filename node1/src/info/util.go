package info

import (
	"math/rand"
	"strings"
	"time"
)

var src *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(n int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	indexBits := 6                       // 6 bits to represent a letter index
	indexMask := int64(1<<indexBits - 1) // All 1-bits, as many as letterIdxBits
	indexMax := 63 / indexBits           // # of letter indices fitting in 63 bits

	sb := strings.Builder{}
	sb.Grow(n)

	for i, cache, remain := n-1, src.Int63(), indexMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), indexMax
		}
		if idx := int(cache & indexMask); idx < len(characters) {
			sb.WriteByte(characters[idx])
			i--
		}
		cache >>= indexBits
		remain--
	}

	return sb.String()
}

func FastRemove(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
