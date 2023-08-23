package components

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestNewField(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	testField1 := NewField("test1")

	expected := Field{
		Name: "test1",
	}

	if expected.Name != testField1.Name {
		t.Fatalf("Entry: expected %v, got %v\n", expected.Name, testField1.Name)
	}
}
