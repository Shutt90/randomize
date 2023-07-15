package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shutt90/password-generator/helpers.go"
)

type storedPassword struct {
	WebsiteName string    `json:"websiteName"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Created     time.Time `json:"created"`
}

var welcomeMessages = []string{
	"Welcome to Randomize Password Manager",
	"",
	"I am constantly looking to improve on this software so please email any suggestions to",
	"Liam.Pugh.009@gmail.com",
	"Any feedback is greatly appreciated",
	"",
	"I hope you enjoy using the product",
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(480, 640))

	welcomeContainer := helpers.CreateTextContainer(welcomeMessages)

	tabs := container.NewVBox(
		container.NewAppTabs(
			container.NewTabItemWithIcon("Home", theme.HomeIcon(), welcomeContainer),
			container.NewTabItemWithIcon("Passwords", theme.ComputerIcon(), widget.NewLabel("Passwords")),
		),
	)

	spacer := widget.NewSeparator()
	spacer.BaseWidget.Resize(fyne.Size{
		Width:  tabs.MinSize().Width,
		Height: tabs.MinSize().Height,
	})

	tabContainer := container.NewHBox(tabs, spacer)

	myWindow.SetContent(tabContainer)

	myWindow.ShowAndRun()

	// fmt.Println("Store or Get?")
	// godotenv.Load()
	// var decision string
	// _, err := fmt.Scanln(&decision)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// if strings.ToLower(decision) == "store" {
	// 	storeCli()
	// 	return
	// }
	// if strings.ToLower(decision) == "get" {
	// 	getCli()
	// 	return
	// }

	// fmt.Println("unknown decision")
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

func storeCli() {
	fmt.Println("Enter website name: ")
	var entry storedPassword
	_, err := fmt.Scanln(&entry.WebsiteName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if os.Getenv("DEFAULT_USERNAME") != "" {
		fmt.Println("(default: ", os.Getenv("DEFAULT_USERNAME"), ")", "Enter username: ")
	}
	_, err = fmt.Scanln(&entry.Username)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Setenv("DEFAULT_USERNAME", entry.Username)
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
	fmt.Println("your password is ", entry.Password)

	ctx := context.Background()

	entry.store(ctx)
}

func getCli() {
	fmt.Println("Enter website name: ")
	var websiteName string
	_, err := fmt.Scanln(&websiteName)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	pw, err := getPassword(websiteName, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(pw)
}

func (sp *storedPassword) store(ctx context.Context) error {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer db.Close()

	conn, err := db.Conn(ctx)
	if err != nil {
		fmt.Println("could not connect to database")
		return nil
	}

	conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS password (websiteName varchar(255), username varchar(255), password varchar(255))")

	query := fmt.Sprintf("INSERT INTO password (websiteName, username, password) VALUES (%v, %v, %v);", &sp.WebsiteName, &sp.Username, &sp.Password)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	fmt.Println("successfully stored password for ", sp.WebsiteName)

	return nil
}

func getPassword(websiteName string, ctx context.Context) (string, error) {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer db.Close()

	conn, err := db.Conn(ctx)
	if err != nil {
		return "", err
	}

	var password string
	err = conn.QueryRowContext(ctx, "SELECT password FROM password WHERE websiteName = ?;", websiteName).Scan(password)
	if err != nil {
		return "", err
	}

	return password, nil
}
