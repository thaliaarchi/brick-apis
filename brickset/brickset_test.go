package brickset

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

var (
	apiKey   = os.Getenv("BRICKSET_API_KEY")
	username = os.Getenv("BRICKSET_USERNAME")
	password = os.Getenv("BRICKSET_PASSWORD")
	userHash = os.Getenv("BRICKSET_USER_HASH")
)

func TestLogin(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	c := NewClient()
	userHash, err := c.Login(apiKey, username, password)
	if err != nil {
		t.Error(err)
	}
	t.Error(userHash)
}

func TestGetSet(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	c := NewClient()
	res, err := c.GetSet(apiKey, userHash, "22667")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(res)
}

func TestGetSets(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	c := NewClient()
	res, err := c.GetSets(apiKey, userHash, "", "", "", "", "", "", "", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(res)
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

func TestEncodeLoginXml(t *testing.T) {
	data, err := xml.Marshal(&loginResponse{Response: "value"})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `<string xmlns="https://brickset.com/api/">test</string>` {
		t.Error("unexpected result", string(data))
	}
}

func TestEncodeGetSetsXml(t *testing.T) {
	data, err := xml.Marshal(GetSetsResponse{Sets: []GetSetResponseItem{{SetID: 1234}}})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "abc" {
		t.Error("unexpected result", string(data))
	}
}

func TestDecodeGetSetsXml(t *testing.T) {
	xmlString := `
	<ArrayOfSets xmlns="https://brickset.com/api/">
	  <sets>
		<setID>22667</setID>
		<number>001</number>
		<numberVariant>1</numberVariant>
		<name>Gears</name>
		<year>1965</year>
		<theme>Samsonite</theme>
		<themeGroup>Vintage themes</themeGroup>
		<subtheme>Basic Set</subtheme>
		<pieces>43</pieces>
		<minifigs />
		<image>true</image>
		<imageFilename>001-1</imageFilename>
		<thumbnailURL>https://images.brickset.com/sets/small/001-1.jpg</thumbnailURL>
		<largeThumbnailURL>https://images.brickset.com/sets/small/001-1.jpg</largeThumbnailURL>
		<imageURL>https://images.brickset.com/sets/images/001-1.jpg</imageURL>
		<bricksetURL>https://brickset.com/sets/001-1</bricksetURL>
		<released>true</released>
		<owned>false</owned>
		<wanted>false</wanted>
		<qtyOwned>0</qtyOwned>
		<userNotes />
		<ACMDataCount>0</ACMDataCount>
		<ownedByTotal>102</ownedByTotal>
		<wantedByTotal>62</wantedByTotal>
		<UKRetailPrice />
		<USRetailPrice>4.95</USRetailPrice>
		<CARetailPrice />
		<EURetailPrice />
		<USDateAddedToSAH />
		<USDateRemovedFromSAH />
		<rating>0</rating>
		<reviewCount>0</reviewCount>
		<packagingType>Box</packagingType>
		<availability>Retail</availability>
		<instructionsCount>0</instructionsCount>
		<additionalImageCount>0</additionalImageCount>
		<ageMin>5</ageMin>
		<ageMax>12</ageMax>
		<height>20.3</height>
		<width>39.5</width>
		<depth>5.1</depth>
		<weight />
		<category>Normal</category>
		<userRating>0</userRating>
		<EAN />
		<UPC />
		<lastUpdated>2018-01-29T10:24:39.983</lastUpdated>
	  </sets>
	</ArrayOfSets>`
	r := &GetSetsResponse{}
	if err := xml.NewDecoder(strings.NewReader(xmlString)).Decode(r); err != nil {
		t.Fatal(err)
	}
	if len(r.Sets) != 1 {
		t.Error("Expected correct string value", len(r.Sets))
	}
}
