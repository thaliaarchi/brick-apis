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
