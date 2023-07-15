package helpers

import (
	"fmt"
	"image/color"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateTextContainer(textArr []string) *fyne.Container {
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

func Randomize(numLetters uint8) string {
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
