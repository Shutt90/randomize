package helpers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Randomize(numLetters uint8) string {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	var password []string
	var i uint8
	for i = 0; i < numLetters; i++ {
		rand := r.Intn(126-33) + 33
		letter := fmt.Sprintf("%c", rand)
		password = append(password, string(letter))
	}

	var sep string
	joinedPassword := strings.Join(password, sep)

	return joinedPassword
}
