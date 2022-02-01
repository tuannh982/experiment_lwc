package commons

import "math/rand"

type Letters string

const (
	Lowercase Letters = "abcdefghijklmnopqrstuvwxyz"
	Uppercase Letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric   Letters = "0123456789"
)

func StringRandom(length int, letters Letters) string {
	lettersLen := len(letters)
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(lettersLen)]
	}
	return string(b)
}
