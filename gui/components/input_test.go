package components

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func TestGetTextBoxes(t *testing.T) {
	testField := Fields{NewField("test"), NewField("test2")}

	res := testField.GetTextBoxes()
	expected := []fyne.CanvasObject{
		canvas.NewText("test", color.White),
		canvas.NewText("test", color.White),
	}

	for i, each := range expected {
		if each != res[i] {
			t.Fatalf("getTextBoxes Failed: expected: %v, got: %v", res[i], each)
		}
	}
}
