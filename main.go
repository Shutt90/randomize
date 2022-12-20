package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	fmt.Println(randomize(16))
}

func randomize(numLetters int) string {
	rand.Seed(time.Now().UnixNano())
	var password []string
	for i := 0; i < numLetters; i++ {
		rand := rand.Intn(126-33) + 33
		letter := fmt.Sprintf("%c", rand)
		password = append(password, string(letter))
	}

	var sep string
	joinedPassword := strings.Join(password, sep)

	return joinedPassword
}
