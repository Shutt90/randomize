package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

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
)

type Fields []Field

type Field struct {
	Name    string
	Entry   *widget.Entry
	Label   *widget.Label
	Textbox canvas.Text
}

func NewField(name string) Field {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(name)
	return Field{
		Name:    name,
		Entry:   entry,
		Label:   widget.NewLabelWithStyle(name, fyne.TextAlignCenter, fyne.TextStyle{}),
		Textbox: *canvas.NewText(name, color.White),
	}
}

func (f Fields) GetTextBoxes() []fyne.CanvasObject {
	var textboxes []fyne.CanvasObject
	for _, field := range f {
		textboxes = append(textboxes, &field.Textbox)
	}

	return textboxes
}

func (f Fields) GetInputs() []fyne.CanvasObject {
	var textboxes []fyne.CanvasObject
	for _, field := range f {
		textboxes = append(textboxes, field.Entry)
	}

	return textboxes
}

func (f Fields) GetInputsWithLabels() []fyne.CanvasObject {
	var labelThenInputOrdered []fyne.CanvasObject
	for _, field := range f {
		labelThenInputOrdered = append(labelThenInputOrdered, field.Label)
		labelThenInputOrdered = append(labelThenInputOrdered, field.Entry)

	}

	return labelThenInputOrdered
}

func (f Fields) MapNamesGetInputs() (map[string]widget.Entry, error) {
	names := make(map[string]widget.Entry)
	for _, field := range f {
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
