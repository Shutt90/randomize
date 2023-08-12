package gui

import (
	"bytes"
	"encoding/json"
	"image/color"
	"net/http"
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
)

type PopupWindow int

type register struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	FirstName     string `json:"firstName"`
	Surname       string `json:"surname"`
	EmailAddress  string `json:"emailAddress"`
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	PostCode      string `json:"postCode"`
}

const (
	Website         = "Website Name"
	Username        = "Username"
	Password        = "Password"
	ConfirmPassword = "Confirm Password"
	FirstName       = "First Name"
	Surname         = "Surname"
	EmailAddress    = "Email Address"
	StreetAddress   = "Street Address"
	City            = "City"
	PostCode        = "Post Code/Zip Code"
	WindowWidth     = 500
	WindowHeight    = 720
)

var regFields = []string{
	Username,
	Password,
	ConfirmPassword,
	FirstName,
	Surname,
	EmailAddress,
	StreetAddress,
	City,
	PostCode,
}

var loginFields = []string{
	Username,
	Password,
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

	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(WindowWidth, WindowHeight))

	switchToRegisterBtn := widget.NewButtonWithIcon("Register", theme.ComputerIcon(), func() {})
	switchToLoginBtn := widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {})

	loginPopup := createLoginMenu(myWindow.Canvas(), switchToRegisterBtn)
	registerPopup := createRegisterMenu(myWindow.Canvas(), switchToLoginBtn)

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

	fields := components.Fields{
		{
			Name:  Website,
			Entry: widget.NewEntry(),
		},
		{
			Name:  Username,
			Entry: widget.NewEntry(),
		},
		{
			Name:  Password,
			Entry: widget.NewEntry(),
		},
	}

	for _, field := range fields {
		field.Textbox = *canvas.NewText(field.Name, color.White)
		field.Textbox.Alignment = fyne.TextAlignCenter
		field.Textbox.TextStyle.Bold = true
		field.Entry.PlaceHolder = field.Name
	}

	pwArr := []fyne.CanvasObject{
		container.NewGridWithColumns(
			len(fields),
			fields.GetTextBoxes()...,
		),
	}
	input := container.NewGridWithColumns(
		len(fields),
		fields.GetInputs()...,
	)

	for _, pass := range passwords {
		passwordContainer := container.NewGridWithColumns(
			len(fields),
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

	storePwButton := widget.NewButton("Store", func() {
		mappedInputsByNames, err := fields.MapNamesGetInputs()
		if err != nil {
			popupForError(myWindow.Canvas(), err.Error())
		}
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
		fields[2].Entry.Text = helpers.Randomize(128)
		fields[2].Entry.Refresh()
	})

	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(fields[2].Entry.Text))
	})

	generateContainer := container.NewGridWithColumns(
		2,
		generatePasswordBtn,
		copyBtn,
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

func createLoginMenu(c fyne.Canvas, btn *widget.Button) *widget.PopUp {
	loginInputArr := components.Fields{}
	items := []fyne.CanvasObject{}

	for _, field := range loginFields {
		loginInputArr = append(loginInputArr, components.NewField(field))
	}

	entries, err := loginInputArr.MapNamesGetInputs()
	if err != nil {
		popupForError(c, err.Error())
	}

	items = append(items, loginInputArr.GetInputsWithLabels()...)

	items = append(
		items,
		widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {
			// make api request when server setup and hide modal
			loginDetails := map[string]string{
				Username: entries["username"].Text,
				Password: entries["password"].Text,
			}

			loginForTransport, err := json.Marshal(&loginDetails)
			if err != nil {
				widget.ShowPopUpAtPosition(
					canvas.NewText("Error while logging in", color.White),
					c,
					fyne.Position{
						X: c.Size().Width/2. - c.Size().Width/2,
						Y: c.Size().Height/2. - c.Size().Height/2,
					},
				)
			}

			http.Post("endpoint", "application/json", bytes.NewBuffer(loginForTransport))
		}),
		btn)

	// TODO: clean this up at earliest possible convience, broken code
	contents := container.NewVBox(items...)

	// Set the desired size for the loginMenu modal
	loginMenuWidth := float32(200.)
	loginMenuHeight := float32(200.)
	loginMenuSize := fyne.NewSize(loginMenuWidth, loginMenuHeight)
	contents.Resize(loginMenuSize)

	loginMenu := widget.NewModalPopUp(contents, c)
	loginMenu.Resize(loginMenuSize) // Set the size of the modal popup

	return loginMenu
}

func createRegisterMenu(c fyne.Canvas, btn *widget.Button) *widget.PopUp {
	entries := []fyne.CanvasObject{}
	regInputArr := components.Fields{}

	for _, regField := range regFields {
		regInputArr = append(regInputArr, components.NewField(regField))
	}

	entries = append(entries, widget.NewButtonWithIcon("Register", theme.DocumentSaveIcon(), func() {
		regFields, err := regInputArr.MapNamesGetInputs()
		if err != nil {
			popupForError(c, err.Error())
		}

		if regFields["password"].Text != regFields["confirmpass"].Text {
			popupForError(c, "passwords do not match")
		}

		registration := register{
			Username:      regFields["username"].Text,
			Password:      regFields["password"].Text,
			FirstName:     regFields["firstname"].Text,
			Surname:       regFields["surname"].Text,
			EmailAddress:  regFields["email"].Text,
			StreetAddress: regFields["street"].Text,
			City:          regFields["city"].Text,
			PostCode:      regFields["postcode"].Text,
		}

		registerForTransport, err := json.Marshal(&registration)
		if err != nil {
			popupForError(c, "error while registering")
		}

		http.Post("endpoint", "application/json", bytes.NewBuffer(registerForTransport))
	}))

	entries = append(entries, btn)

	contents := container.NewVBox(
		entries...,
	)
	// Set the desired size for the loginMenu modal
	registerMenuWidth := float32(200.)
	registerMenuHeight := float32(200.)
	registerMenuSize := fyne.NewSize(registerMenuWidth, registerMenuHeight)
	contents.Resize(registerMenuSize)

	registerMenu := widget.NewModalPopUp(contents, c)
	registerMenu.Resize(registerMenuSize) // Set the size of the modal popup

	return registerMenu
}

func popupForError(c fyne.Canvas, msg string) {
	widget.ShowPopUpAtPosition(
		canvas.NewText(msg, color.White),
		c,
		fyne.Position{
			X: c.Size().Width/2. - c.Size().Width/2,
			Y: c.Size().Height/2. - c.Size().Height/2,
		},
	)
}
