package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/shutt90/password-generator/helpers.go"
)

const (
	Register PopupWindow = iota
	Login
)

type PopupWindow int

type fields []field

type field struct {
	Name    string
	Entry   *widget.Entry
	Label   *widget.Label
	Textbox canvas.Text
}

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

var welcomeMessages = []string{
	"Welcome to Randomize Password Manager",
	"",
	"I am constantly looking to improve on this software so please email any suggestions to",
	"Liam.Pugh.009@gmail.com",
	"Any feedback is greatly appreciated",
	"",
	"I hope you enjoy using the product",
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

	fields := fields{
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
			fields.getTextBoxes()...,
		),
	}
	input := container.NewGridWithColumns(
		len(fields),
		fields.getInputs()...,
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
		mappedInputsByNames, err := fields.mapNamesGetInputs()
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

func createLoginMenu(c fyne.Canvas, b *widget.Button) *widget.PopUp {
	loginFields := fields{
		{
			Name:  Username,
			Label: widget.NewLabelWithStyle(Username, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  Password,
			Label: widget.NewLabelWithStyle(Password, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
	}

	contents := container.NewVBox(
		loginFields[0].Label,
		loginFields[0].Entry,
		loginFields[1].Label,
		loginFields[1].Entry,
		widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {
			// make api request when server setup and hide modal
			loginDetails := map[string]string{
				Username: loginFields[0].Entry.Text,
				Password: loginFields[1].Entry.Text,
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
		b,
	)

	// Set the desired size for the loginMenu modal
	loginMenuWidth := float32(200.)
	loginMenuHeight := float32(200.)
	loginMenuSize := fyne.NewSize(loginMenuWidth, loginMenuHeight)
	contents.Resize(loginMenuSize)

	loginMenu := widget.NewModalPopUp(contents, c)
	loginMenu.Resize(loginMenuSize) // Set the size of the modal popup

	return loginMenu
}

func createRegisterMenu(c fyne.Canvas, b *widget.Button) *widget.PopUp {
	entries := []fyne.CanvasObject{}
	regInputs := fields{
		{
			Name:  Username,
			Label: widget.NewLabelWithStyle(Username, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  FirstName,
			Label: widget.NewLabelWithStyle(FirstName, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  Surname,
			Label: widget.NewLabelWithStyle(Surname, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  Password,
			Label: widget.NewLabelWithStyle(Password, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  ConfirmPassword,
			Label: widget.NewLabelWithStyle(ConfirmPassword, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  EmailAddress,
			Label: widget.NewLabelWithStyle(EmailAddress, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  StreetAddress,
			Label: widget.NewLabelWithStyle(StreetAddress, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  City,
			Label: widget.NewLabelWithStyle(City, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
		{
			Name:  PostCode,
			Label: widget.NewLabelWithStyle(PostCode, fyne.TextAlignCenter, fyne.TextStyle{}),
			Entry: widget.NewEntry(),
		},
	}

	for _, input := range regInputs {
		entries = append(entries, input.Label)
		entries = append(entries, input.Entry)
	}

	entries = append(entries, widget.NewButtonWithIcon("Register", theme.DocumentSaveIcon(), func() {
		regFields, err := regInputs.mapNamesGetInputs()
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

	entries = append(entries, b)

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

func (f fields) getTextBoxes() []fyne.CanvasObject {
	var textboxes []fyne.CanvasObject
	for _, field := range f {
		textboxes = append(textboxes, &field.Textbox)
	}

	return textboxes
}

func (f fields) getInputs() []fyne.CanvasObject {
	var textboxes []fyne.CanvasObject
	for _, field := range f {
		textboxes = append(textboxes, field.Entry)
	}

	return textboxes
}

func (f fields) mapNamesGetInputs() (map[string]widget.Entry, error) {
	names := make(map[string]widget.Entry)
	for _, field := range f {
		if field.Entry.Text == "" {
			return map[string]widget.Entry{}, fmt.Errorf("all fields must be filled in")
		}
		switch field.Name {
		case Website:
			names["website"] = *field.Entry
		case Username:
			names["username"] = *field.Entry
		case Password:
			names["password"] = *field.Entry
		case ConfirmPassword:
			names["confirmpass"] = *field.Entry
		case FirstName:
			names["firstname"] = *field.Entry
		case Surname:
			names["surname"] = *field.Entry
		case EmailAddress:
			names["email"] = *field.Entry
		case StreetAddress:
			names["street"] = *field.Entry
		case City:
			names["city"] = *field.Entry
		case PostCode:
			names["postcode"] = *field.Entry
		}
	}

	return names, nil
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
