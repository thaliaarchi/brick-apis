package bricklinkstore

import (
	"testing"

	"github.com/andrewarchi/bricklink-buy/credentials"
)

func TestColors(t *testing.T) {
	cred, err := credentials.Read("../credentials.json")
	if err != nil {
		t.Fatal(err)
	}
	bl, err := NewClient(cred.BrickLinkStore)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := bl.GetColorList()
	if err != nil {
		t.Fatal(err)
	}
	if len(colors) == 0 {
		t.Error("expected colors response, but got none")
	}
	for _, c := range colors {
		color, err := bl.GetColor(c.ColorID)
		if err != nil {
			t.Error(err)
		}
		if c != *color {
			t.Errorf("expected %v, but got %v", c, color)
		}
	}
}
