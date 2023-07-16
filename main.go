package main

import (
	"context"
	"database/sql"
	"image/color"
	"log"
	"os"

	cockroachDB "github.com/shutt90/password-generator/db"
	"github.com/shutt90/password-generator/helpers.go"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	err := godotenv.Load()
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer db.Close()

	conn, err := db.Conn(ctx)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	cc := cockroachDB.NewCockroachClient(ctx, db)

	passwords, err := cc.GetAllPasswords()
	if err != nil {
		panic(err)
	}
	myApp := app.New()

	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(480, 640))

	welcomeContainer := helpers.CreateTextContainer(welcomeMessages)

	storePwButton := widget.NewButton("Store", func() {
		log.Println("tapped")
	})

	storePwButton.Alignment = widget.ButtonAlign(fyne.TextAlignCenter)
	pwArr := []fyne.CanvasObject{
		container.NewGridWithColumns(
			3,
			canvas.NewText("Username", color.White),
			canvas.NewText("WebsiteName", color.White),
			canvas.NewText("Password", color.White),
		),
	}
	for _, pass := range passwords {
		passwordContainer := container.NewGridWithColumns(
			3,
			canvas.NewText(pass.Username, color.White),
			canvas.NewText(pass.WebsiteName, color.White),
			canvas.NewText(pass.Password, color.White),
		)

		pwArr = append(pwArr, passwordContainer)
	}

	pwContainer := container.NewPadded(
		container.NewVBox(
			container.NewVBox(
				container.NewVBox(pwArr...),
			),
		),
	)

	tabs := container.NewVBox(
		container.NewAppTabs(
			container.NewTabItemWithIcon("Home", theme.HomeIcon(), welcomeContainer),
			container.NewTabItemWithIcon("Passwords", theme.ComputerIcon(),
				container.NewVBox(pwContainer, storePwButton),
			),
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
