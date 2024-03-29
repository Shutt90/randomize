package components

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type register struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	FirstName     string `json:"firstName"`
	Surname       string `json:"surname"`
	Email         string `json:"email"`
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	PostCode      string `json:"postCode"`
}

func CreatePopup(c fyne.Canvas, btn *widget.Button, fields Fields, toRegister bool) (*widget.PopUp, error) {
	items := []fyne.CanvasObject{}
	entries, err := fields.MapNamesGetInputs()
	if err != nil {
		return nil, err
	}

	items = append(items, fields.GetInputsWithLabels()...)

	items = append(
		items,
		widget.NewButtonWithIcon("Login", theme.LoginIcon(), func() {
			// make api request when server setup and hide modal
			sliceOfEntries := []string{}
			for _, entry := range entries {
				sliceOfEntries = append(sliceOfEntries, entry.Text)
			}

			dataForTransport := &bytes.Buffer{}

			err := gob.NewEncoder(dataForTransport).Encode(sliceOfEntries)
			if err != nil {
				fmt.Println(err)
				return
			}

			// var res *http.Response

			if toRegister {
				// res, err = http.Post("endpoint", "application/json", dataForTransport)
				// if err != nil {
				// 	return
				// }
			} else {
				// res, err = http.Post("endpoint", "application/json", dataForTransport)
				// if err != nil {
				// 	return
				// }
			}

			// if res.StatusCode == 200 {

			// }
		}),
		btn,
	)

	contents := container.NewVBox(items...)

	// Set the desired size for the loginMenu modal
	popupMenuWidth := float32(200.)
	popupMenuHeight := float32(200.)
	popupMenuSize := fyne.NewSize(popupMenuWidth, popupMenuHeight)
	contents.Resize(popupMenuSize)

	popupMenu := widget.NewModalPopUp(contents, c)
	popupMenu.Resize(popupMenuSize) // Set the size of the modal popup

	return popupMenu, nil
}
