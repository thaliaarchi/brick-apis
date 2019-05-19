package bricklinkstore

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrewarchi/bricklink-buy/credentials"
	"github.com/mrjones/oauth"
)

const base = "https://api.bricklink.com/api/store/v1"

// Client connects to the BrickLink store API.
type Client struct {
	client *http.Client
}

// NewClient constructs a client for the BrickLink store API.
func NewClient(cred *credentials.BrickLinkStore) (*Client, error) {
	consumer := oauth.NewConsumer(cred.ConsumerKey, cred.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cred.Token, Secret: cred.TokenSecret}
	client, err := consumer.MakeHttpClient(accessToken)
	return &Client{client}, err
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(base + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func checkMeta(m meta) error {
	if m.Code/100 != 2 {
		return fmt.Errorf("status code not OK: %d %s (%s)", m.Code, m.Message, m.Description)
	}
	return nil
}

type meta struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Code        int    `json:"code"`
}
