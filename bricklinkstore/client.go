package bricklinkstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mrjones/oauth"
)

const base = "https://api.bricklink.com/api/store/v1"

// Client connects to the BrickLink store API.
type Client struct {
	client *http.Client
}

// NewClient constructs a client for the BrickLink store API.
func NewClient(consumerKey, consumerSecret, token, tokenSecret string) (*Client, error) {
	consumer := oauth.NewConsumer(consumerKey, consumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: token, Secret: tokenSecret}
	client, err := consumer.MakeHttpClient(accessToken)
	return &Client{client}, err
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(base + url)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status %s", resp.Status)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
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
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status %s", resp.Status)
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
		decoder.DisallowUnknownFields()
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

// 2xx is successful: http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki/Error_Handling
func checkMeta(m meta) error {
	if m.Code/100 != 2 {
		return fmt.Errorf("status %d %s: %s", m.Code, m.Message, m.Description)
	}
	return nil
}

type meta struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Code        int64  `json:"code"`
}
