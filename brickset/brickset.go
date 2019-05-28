package brickset

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

const (
	url      = "https://brickset.com/api"
	endpoint = url + "/v2.asmx"
)

func getFormEncodedEndpoint(methodName string) string {
	return fmt.Sprintf("%s/%s", endpoint, methodName)
}

// Client enables the ability to make requests to the Brickset API
type Client struct {
	c *http.Client
}

// NewClient creates a new Brickset client
func NewClient() *Client {
	return &Client{&http.Client{}}
}

func (c *Client) makeRequest(method string, body string, result interface{}) error {
	req, _ := http.NewRequest("POST", getFormEncodedEndpoint(method), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return xml.NewDecoder(resp.Body).Decode(result)
}

// Login is used to login to the Brickset API
func (c Client) Login(apiKey, username, password string) (string, error) {
	body := fmt.Sprintf(`apiKey=%s&username=%s&password=%s`, apiKey, username, password)
	r := &loginResponse{}
	if err := c.makeRequest("login", body, r); err != nil {
		return "", err
	}
	return r.Response, nil
}

type loginResponse struct {
	XMLName  xml.Name `xml:"https://brickset.com/api/ string"`
	Response string   `xml:",chardata"`
}

// GetSet is used to get the details for a single set
func (c *Client) GetSet(apiKey, userHash, setID string) (*GetSetsResponse, error) {
	body := fmt.Sprintf("apiKey=%s&userHash=%s&SetID=%s", apiKey, userHash, setID)
	r := &GetSetsResponse{}
	if err := c.makeRequest("getSet", body, r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetSets is used to get a list of sets from the Brickset API
func (c *Client) GetSets(apiKey, userHash, query, theme, subTheme, setNumber, year, owned, wanted, orderBy, pageSize, pageNumber, userName string) (*GetSetsResponse, error) {
	body := fmt.Sprintf("apiKey=%s&userHash=%s&query=%s&theme=%s&subtheme=%s&setNumber=%s&year=%s&owned=%s&wanted=%s&orderBy=%s&pageSize=%s&pageNumber=%s&userName=%s", apiKey, userHash, query, theme, subTheme, setNumber, year, owned, wanted, orderBy, pageSize, pageNumber, userName)
	r := &GetSetsResponse{}
	if err := c.makeRequest("getSets", body, r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetSetsResponse is the data returned from getSets
type GetSetsResponse struct {
	XMLName xml.Name             `xml:"https://brickset.com/api/ ArrayOfSets"`
	Sets    []GetSetResponseItem `xml:"sets"`
}

// GetSetResponseItem is a single set returned by getSets
type GetSetResponseItem struct {
	XMLName              xml.Name `xml:"sets"`
	SetID                int      `xml:"setID"`
	Number               string   `xml:"number"`
	NumberVariant        int      `xml:"numberVariant"`
	Name                 string   `xml:"name"`
	Year                 string   `xml:"year"`
	Theme                string   `xml:"theme"`
	ThemeGroup           string   `xml:"themeGroup"`
	Subtheme             string   `xml:"subtheme"`
	Pieces               string   `xml:"pieces"`
	Minifigs             string   `xml:"minifigs"`
	Image                bool     `xml:"image"`
	ImageFilename        string   `xml:"imageFilename"`
	ThumbnailURL         string   `xml:"thumbnailURL"`
	LargeThumbnailURL    string   `xml:"largeThumbnailURL"`
	ImageURL             string   `xml:"imageURL"`
	BricksetURL          string   `xml:"bricksetURL"`
	Released             bool     `xml:"released"`
	Owned                bool     `xml:"owned"`
	Wanted               bool     `xml:"wanted"`
	QtyOwned             int      `xml:"qtyOwned"`
	UserNotes            string   `xml:"userNotes"`
	ACMDataCount         int      `xml:"ACMDataCount"`
	OwnedByTotal         int      `xml:"ownedByTotal"`
	WantedByTotal        int      `xml:"wantedByTotal"`
	UKRetailPrice        string   `xml:"UKRetailPrice"`
	USRetailPrice        string   `xml:"USRetailPrice"`
	CARetailPrice        string   `xml:"CARetailPrice"`
	EURetailPrice        string   `xml:"EURetailPrice"`
	USDateAddedToSAH     string   `xml:"USDateAddedToSAH"`
	USDateRemovedFromSAH string   `xml:"USDateRemovedFromSAH"`
	Rating               float32  `xml:"rating"`
	ReviewCount          int      `xml:"reviewCount"`
	PackagingType        string   `xml:"packagingType"`
	Availability         string   `xml:"availability"`
	InstructionsCount    int      `xml:"instructionsCount"`
	AdditionalImageCount int      `xml:"additionalImageCount"`
	AgeMin               string   `xml:"ageMin"`
	AgeMax               string   `xml:"ageMax"`
	Height               string   `xml:"height"`
	Width                string   `xml:"width"`
	Depth                string   `xml:"depth"`
	Weight               string   `xml:"weight"`
	Category             string   `xml:"category"`
	Notes                string   `xml:"notes"`
	UserRating           string   `xml:"userRating"`
	Tags                 string   `xml:"tags"`
	EAN                  string   `xml:"EAN"`
	UPC                  string   `xml:"UPC"`
	Description          string   `xml:"description"`
	LastUpdated          string   `xml:"lastUpdated"`
}
