package gui

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	cockroachDB "github.com/shutt90/password-generator/db"
	"github.com/shutt90/password-generator/helpers.go"
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

func MainWindow(db *cockroachDB.CockroachClient, passwords []cockroachDB.StoredPassword) {
	myApp := app.New()

	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(480, 640))

	infoContainer := createTextContainer(welcomeMessages)

	var fieldArr []fyne.CanvasObject
	inputArr := make([]*widget.Entry, len(fieldArr))

	for _, field := range fields {
		headerField := canvas.NewText(field, color.White)
		headerField.Alignment = fyne.TextAlignCenter
		headerField.TextStyle.Bold = true
		fieldArr = append(fieldArr, headerField)
		inputField := widget.NewEntry()
		inputField.PlaceHolder = field
		inputArr = append(inputArr, inputField)
	}

	pwArr := []fyne.CanvasObject{
		container.NewGridWithColumns(
			len(fieldArr),
			fieldArr...,
		),
	}
	input := container.NewGridWithColumns(
		len(inputArr),
		fyne.CanvasObject(inputArr[0]),
		fyne.CanvasObject(inputArr[1]),
		fyne.CanvasObject(inputArr[2]),
	)
	for _, pass := range passwords {
		passwordContainer := container.NewGridWithColumns(
			len(fieldArr),
			canvas.NewText(pass.WebsiteName, color.White),
			canvas.NewText(pass.Username, color.White),
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

	storePwButton := widget.NewButton("Store", func() {
		input := cockroachDB.StoredPassword{
			WebsiteName: inputArr[0].Text,
			Username:    inputArr[1].Text,
			Password:    inputArr[2].Text,
		}

		err := db.Store(input)
		if err != nil {
			newPaddedH := container.NewHBox(container.NewPadded(canvas.NewText(err.Error(), color.White)))
			widget.ShowPopUpAtPosition(
				newPaddedH,
				myWindow.Canvas(),
				fyne.Position{
					X: myWindow.Canvas().Size().Width/2. - newPaddedH.Size().Width/2,
					Y: myWindow.Canvas().Size().Height/2. - newPaddedH.Size().Height/2,
				},
			)
		}

		//TODO: Add way to refresh passwords on submit
		passwords, err = db.GetAllPasswords()
		if err != nil {
			panic(err)
		}
	})

	passwordOutput := widget.NewEntry()
	passwordOutput.Disable()
	generatePasswordBtn := widget.NewButton("Generate Password", func() {
		randomize := helpers.Randomize(128)
		inputArr[2].Text = randomize
		inputArr[2].Refresh()
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

func createTextContainer(textArr []string) *fyne.Container {
	var canvObj []fyne.CanvasObject

	for _, text := range textArr {
		if text == "Liam.Pugh.009@gmail.com" {
			url := &url.URL{
				Scheme: "mailto",
				Path:   text,
			}
			hyperlink := widget.NewHyperlink(text, url)
			hyperlink.Alignment = fyne.TextAlignCenter
			canvObj = append(canvObj, hyperlink)
			continue
		}
		textObj := canvas.NewText(text, color.White)
		textObj.Alignment = fyne.TextAlignCenter
		canvObj = append(canvObj, textObj)
	}

	textContainer := container.NewPadded(
		container.NewVBox(canvObj...),
	)

	return textContainer
}
