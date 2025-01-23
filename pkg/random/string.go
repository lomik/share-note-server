package random

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func String(dict string, size int) string {
	var out strings.Builder
	for i := 0; i < size; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(dict))))
		if err != nil {
			panic(err)
		}
		out.WriteByte(dict[int(nBig.Int64())])
	}
	return out.String()
}
