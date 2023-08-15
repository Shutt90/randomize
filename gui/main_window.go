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
	"golang.design/x/clipboard"

	cockroachDB "github.com/shutt90/password-generator/db"
	"github.com/shutt90/password-generator/gui/components"
	"github.com/shutt90/password-generator/helpers.go"
)

const (
	Register PopupWindow = iota
	Login

	WindowWidth  = 500
	WindowHeight = 720
)

type PopupWindow int

type Canvas fyne.Canvas

var regFieldNames = []string{
	components.Username,
	components.Password,
	components.ConfirmPassword,
	components.FirstName,
	components.Surname,
	components.EmailAddress,
	components.StreetAddress,
	components.City,
	components.PostCode,
}

var loginFieldNames = []string{
	components.Username,
	components.Password,
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

func MainWindow(db *cockroachDB.CockroachClient, passwords []cockroachDB.StoredPassword) {
	myApp := app.New()
	verified := make(chan bool, 1)

	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(WindowWidth, WindowHeight))

	mainCanvas := myWindow.Canvas()

	switchToRegisterBtn := widget.NewButtonWithIcon("Register", theme.ComputerIcon(), func() {})
	switchToLoginBtn := widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {})

	loginFields := components.Fields{}
	registerFields := components.Fields{}
	for _, fieldName := range loginFieldNames {
		loginFields = append(loginFields, components.NewField(fieldName))
	}
	for _, fieldName := range regFieldNames {
		registerFields = append(registerFields, components.NewField(fieldName))
	}

	loginPopup, err := components.CreatePopup(mainCanvas, switchToRegisterBtn, loginFields, false, verified)
	if err != nil {
		popupForError(myWindow.Canvas(), err.Error())
		return
	}
	registerPopup, err := components.CreatePopup(mainCanvas, switchToRegisterBtn, registerFields, true, verified)
	if err != nil {
		popupForError(myWindow.Canvas(), err.Error())
		return
	}

	switchToRegisterBtn.OnTapped = func() {
		loginPopup.Hide()
		registerPopup.Show()
	}
	switchToLoginBtn.OnTapped = func() {
		registerPopup.Hide()
		loginPopup.Show()
	}

	loginPopup.Show()

	infoContainer := createTextContainer(welcomeMessages)

	newPasswordInputs := components.Fields{}
	for _, inputField := range []string{components.Website, components.Username, components.Password} {
		newField := components.NewField(inputField)
		newField.Textbox.TextStyle.Bold = true
		newPasswordInputs = append(newPasswordInputs, newField)
	}
	pwArr := []fyne.CanvasObject{
		container.NewGridWithColumns(
			len(newPasswordInputs),
			newPasswordInputs.GetTextBoxes()...,
		),
	}
	input := container.NewGridWithColumns(
		len(newPasswordInputs),
		newPasswordInputs.GetInputs()...,
	)

	for _, pass := range passwords {
		passwordContainer := container.NewGridWithColumns(
			len(newPasswordInputs),
			canvas.NewText(pass.WebsiteName, color.White),
			canvas.NewText(pass.Username, color.White),
			canvas.NewText(pass.Password, color.White),
		)

		pwArr = append(pwArr, passwordContainer)
	}

	pwContainer := container.NewPadded(
		container.NewVBox(
			container.NewVBox(pwArr...),
		),
	)

	mappedInputsByNames, err := newPasswordInputs.MapNamesGetInputs()
	if err != nil {
		popupForError(myWindow.Canvas(), err.Error())
	}

	storePwButton := widget.NewButton("Store", func() {

		input := cockroachDB.StoredPassword{
			WebsiteName: mappedInputsByNames["website"].Text,
			Username:    mappedInputsByNames["username"].Text,
			Password:    mappedInputsByNames["password"].Text,
		}

		err = db.Store(input)
		if err != nil {
			newPaddedH := container.NewHBox(container.NewPadded(canvas.NewText(err.Error(), color.White)))
			widget.ShowPopUpAtPosition(
				newPaddedH,
				myWindow.Canvas(),
				fyne.Position{
					X: mainCanvas.Size().Width/2. - newPaddedH.Size().Width/2,
					Y: mainCanvas.Size().Height/2. - newPaddedH.Size().Height/2,
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
		pwInput := mappedInputsByNames["password"]

		pwInput.Text = helpers.Randomize(128)
		pwInput.Refresh()
	})

	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(mappedInputsByNames["password"].Text))
	})

	generateContainer := container.NewGridWithColumns(
		2,
		generatePasswordBtn,
		copyBtn,
	)

	infoTab := container.NewTabItemWithIcon("Info", theme.InfoIcon(), infoContainer)

	tabs := container.NewAppTabs(
		infoTab,
	)

	pwTab := container.NewTabItemWithIcon("Passwords", theme.ComputerIcon(),
		container.NewVBox(
			pwContainer,
			generateContainer,
			container.NewVBox(
				input,
				storePwButton,
			),
		),
	)

	tabContainer := container.NewVBox(tabs)

	myWindow.SetContent(tabContainer)

	myWindow.ShowAndRun()

	if <-verified {
		tabs = container.NewAppTabs(
			infoTab,
			pwTab,
		)

		loginPopup.Hide()
		registerPopup.Hide()

		tabContainer = container.NewVBox(tabs)
		myWindow.SetContent(tabContainer)
		myWindow.Canvas().Refresh(myWindow.Content())
	}
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

func popupForError(c Canvas, msg string) {
	widget.ShowPopUpAtPosition(
		canvas.NewText(msg, color.White),
		c,
		fyne.Position{
			X: c.Size().Width/2. - c.Size().Width/2,
			Y: c.Size().Height/2. - c.Size().Height/2,
		},
	)
}
