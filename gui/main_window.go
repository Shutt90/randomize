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
	"github.com/shutt90/password-generator/helpers.go"
)

type fields []field

type field struct {
	Name    string
	Entry   *widget.Entry
	Textbox canvas.Text
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
	Website  = "Website Name"
	Username = "Username"
	Password = "Password"
)

func MainWindow(db *cockroachDB.CockroachClient, passwords []cockroachDB.StoredPassword) {
	myApp := app.New()

	myWindow := myApp.NewWindow("Randomize Password Manager")
	myWindow.Resize(fyne.NewSize(480, 640))

	//TODO: if no auth token is provided
	createLoginMenu(myWindow.Canvas())

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
		mappedInputsByNames := fields.mapNamesGetInputs()
		input := cockroachDB.StoredPassword{
			WebsiteName: mappedInputsByNames["website"].Text,
			Username:    mappedInputsByNames["username"].Text,
			Password:    mappedInputsByNames["password"].Text,
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

func createLoginMenu(c fyne.Canvas) {
	usernameLabel := widget.NewLabel("Username")
	usernameLabel.Alignment = fyne.TextAlignCenter
	passwordLabel := widget.NewLabel("Password")
	passwordLabel.Alignment = fyne.TextAlignCenter
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewEntry()

	contents := container.NewVBox(
		usernameLabel,
		usernameEntry,
		passwordLabel,
		passwordEntry,
		widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {
			// make api request when server setup and hide modal
			loginDetails := map[string]string{
				Username: usernameEntry.Text,
				Password: passwordEntry.Text,
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
	)
	// Set the desired size for the loginMenu modal
	loginMenuWidth := float32(200.)
	loginMenuHeight := float32(200.)
	loginMenuSize := fyne.NewSize(loginMenuWidth, loginMenuHeight)
	contents.Resize(loginMenuSize)

	loginMenu := widget.NewModalPopUp(contents, c)
	loginMenu.Resize(loginMenuSize) // Set the size of the modal popup

	loginMenu.Show()
}

func createRegisterMenu(c fyne.Canvas) {
	entries := []fyne.CanvasObject{}
	regInputs := []string{
		"Username",
		"Password",
		"First Name",
		"Surname",
		"Confirm Password",
		"Email Address",
		"Street Address",
		"City",
		"Postal Code/Zip Code",
	}

	for _, input := range regInputs {
		iterationInput := widget.NewLabel(input)
		iterationInput.Alignment = fyne.TextAlignCenter
		entries = append(entries, iterationInput)
		entries = append(entries, widget.NewEntry())
	}

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

	registerMenu.Show()

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

func (f fields) mapNamesGetInputs() map[string]widget.Entry {
	names := make(map[string]widget.Entry)
	for _, field := range f {
		switch field.Name {
		case Website:
			names["website"] = *field.Entry
		case Username:
			names["username"] = *field.Entry
		case Password:
			names["password"] = *field.Entry
		}
	}

	return names
}
