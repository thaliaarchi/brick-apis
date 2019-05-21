package brickset

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	apiKey := os.Getenv("BRICKSET_APIKEY")
	username := os.Getenv("BRICKSET_USERNAME")
	password := os.Getenv("BRICKSET_PASSWORD")
	c := NewClient()
	if _, err := c.Login(apiKey, username, password); err != nil {
		t.Error(err)
	}
}

func TestDecodeLoginResponseXML(t *testing.T) {
	xmlString := `<string xmlns="https://brickset.com/api/">test</string>`
	r := &loginResponse{}
	if err := xml.NewDecoder(strings.NewReader(xmlString)).Decode(r); err != nil {
		t.Fatal(err)
	}
	if r.Response != "test" {
		t.Error("Expected correct string value", r.Response)
	}
}

func TestEncodeXml(t *testing.T) {
	data, err := xml.Marshal(&loginResponse{Response: "value"})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `<string xmlns="https://brickset.com/api/">test</string>` {
		t.Error("unexpected result")
	}
}
