package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type storedPassword struct {
	WebsiteName string    `json:"websiteName"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

func main() {
	fmt.Println("Enter website name: ")
	var entry storedPassword
	_, err := fmt.Scanln(&entry.WebsiteName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Enter username: ")
	_, err = fmt.Scanln(&entry.Username)
	if err != nil {
		fmt.Println(err)
		return
	}

	// scanner.Scan(entry.WebsiteName)
	fmt.Println("Number of characters")
	var chars string
	n, err := fmt.Scanln(&chars)
	if err != nil || n > 3 {
		fmt.Println("Error: can't be higher than 255 and ", err)
		return
	}
	convertedNum, err := strconv.Atoi(chars)
	if err != nil {
		fmt.Println(err)
		return
	}
	if convertedNum > 255 || convertedNum < 8 {
		fmt.Println("Error: can't be higher than 255 or lower than 8")
		return
	}

	entry.Password = randomize(uint8(convertedNum))
}

func randomize(numLetters uint8) string {
	rand.Seed(time.Now().UnixNano())
	var password []string
	var i uint8
	for i = 0; i < numLetters; i++ {
		rand := rand.Intn(126-33) + 33
		letter := fmt.Sprintf("%c", rand)
		password = append(password, string(letter))
	}

	var sep string
	joinedPassword := strings.Join(password, sep)

	return joinedPassword
}
