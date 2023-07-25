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

var fields = []string{
	"Website",
	"Username",
	"Password",
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
	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS password (websiteName varchar(255), username varchar(255), password varchar(255))")

	defer conn.Close()

	cc := cockroachDB.NewCockroachClient(ctx, db)

	passwords, err := cc.GetAllPasswords()
	if err != nil {
		panic(err)
	}
	myApp := app.New()

	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(480, 640))

	infoContainer := helpers.CreateTextContainer(welcomeMessages)

	usernameHeader := canvas.NewText("Username", color.White)
	passwordHeader := canvas.NewText("Password", color.White)
	websiteNameHeader := canvas.NewText("Website Name", color.White)

	usernameHeader.TextStyle.Bold = true
	passwordHeader.TextStyle.Bold = true
	websiteNameHeader.TextStyle.Bold = true

	pwArr := []fyne.CanvasObject{
		container.NewGridWithColumns(
			3,
			websiteNameHeader,
			usernameHeader,
			passwordHeader,
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

	websiteField := widget.NewEntry()
	usernameField := widget.NewEntry()
	passwordField := widget.NewEntry()
	websiteField.PlaceHolder = "Website Name"
	usernameField.PlaceHolder = "Login Username"
	passwordField.PlaceHolder = "Password"

	input := container.NewGridWithColumns(
		3,
		websiteField,
		usernameField,
		passwordField,
	)

	storePwButton := widget.NewButton("Store", func() {
		input := cockroachDB.StoredPassword{
			WebsiteName: websiteField.Text,
			Username:    usernameField.Text,
			Password:    passwordField.Text,
		}

		err = cc.Store(input)
		if err != nil {
			widget.ShowPopUpAtPosition(container.NewHBox(canvas.NewText(err.Error(), color.RGBA{255, 0, 0, 1})), myWindow.Canvas(), fyne.Position{X: myWindow.Canvas().Size().Width / 2., Y: myWindow.Canvas().Size().Height / 2.})
		}

		//TODO: Add way to refresh passwords on submit
		passwords, err = cc.GetAllPasswords()
		if err != nil {
			panic(err)
		}
	})

	passwordOutput := widget.NewEntry()
	passwordOutput.Disable()
	generatePasswordBtn := widget.NewButton("Generate Password", func() {
		randomize := helpers.Randomize(128)
		passwordField.Text = randomize
		widget.ShowPopUpAtPosition(container.NewHBox(canvas.NewText("hi", color.White)), myWindow.Canvas(), fyne.Position{X: myWindow.Canvas().Size().Width / 2., Y: myWindow.Canvas().Size().Height / 2.})
	})

	generateContainer := container.NewGridWithColumns(
		1,
		generatePasswordBtn,
	)

	tabs := container.NewVBox(
		container.NewAppTabs(
			container.NewTabItemWithIcon("Passwords", theme.ComputerIcon(),
				container.NewVBox(
					pwContainer,
					generateContainer,
					container.NewVBox(
						input,
						storePwButton,
					),
				),
			),
			container.NewTabItemWithIcon("Info", theme.InfoIcon(), infoContainer),
		),
	)

	spacer := widget.NewSeparator()
	spacer.BaseWidget.Resize(fyne.Size{
		Width:  tabs.MinSize().Width,
		Height: tabs.MinSize().Height,
	})

	tabContainer := container.NewVBox(tabs)

	myWindow.SetContent(tabContainer)

	myWindow.ShowAndRun()
}
