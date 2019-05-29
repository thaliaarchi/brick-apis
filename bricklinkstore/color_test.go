package bricklinkstore

import (
	"os"
	"testing"
)

var (
	consumerKey    = os.Getenv("BRICKLINK_STORE_CONSUMER_KEY")
	consumerSecret = os.Getenv("BRICKLINK_STORE_CONSUMER_SECRET")
	token          = os.Getenv("BRICKLINK_STORE_TOKEN")
	tokenSecret    = os.Getenv("BRICKLINK_STORE_TOKEN_SECRET")
)

func TestColors(t *testing.T) {
	bl, err := NewClient(consumerKey, consumerSecret, token, tokenSecret)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := bl.GetColors()
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
