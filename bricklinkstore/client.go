package bricklinkstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

func (c *Client) doGetAndSave(url string, v interface{}, filename string) error {
	file, err := os.Create("../data/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := c.client.Get(base + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	pr, pw := io.Pipe()
	tr := io.TeeReader(resp.Body, pw)

	done := make(chan bool)
	errs := make(chan error)
	defer close(done)

	go func() {
		defer pw.Close()
		if _, err := io.Copy(file, tr); err != nil {
			errs <- err
		}
		done <- true
	}()

	go func() {
		decoder := json.NewDecoder(pr)
		if err := decoder.Decode(&v); err != nil {
			errs <- err
		}
		done <- true
	}()

	<-done
	<-done
	close(errs)
	err = nil
	for e := range errs {
		if err == nil {
			err = e
		} else {
			err = fmt.Errorf("%s\n%s", err, e)
		}
	}
	return err
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
