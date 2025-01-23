package random

import "strings"

const TokenChars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Token(size int) string {
	return String(TokenChars, size)
}

func IsToken(size int, value string) bool {
	return len(value) == size && strings.Trim(value, TokenChars) == ""
}
