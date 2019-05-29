package bricklinkuser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	client *http.Client
}

func NewClient() (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	return &Client{&http.Client{Jar: jar}}, nil
}

func (c *Client) Login(username, password string) error {
	return c.loginAndOut(url.Values{
		"userid":          {username},
		"password":        {password},
		"keepme_loggedin": {"true"},
	})
}

func (c *Client) Logout() error {
	return c.loginAndOut(url.Values{
		"do_logout": {"true"},
	})
}

// loginReturn is returned by the Login method
type loginReturn struct {
	ReturnCode     int    `json:"returnCode"`
	ReturnMessage  string `json:"returnMessage"`
	ErrorTicket    int    `json:"errorTicket"`
	ProcessingTime int    `json:"procssingTime"`
}

func (c *Client) loginAndOut(formValues url.Values) error {
	url := fmt.Sprintf("https://%s/ajax/renovate/loginandout.ajax", getHost("www"))
	resp, err := c.client.PostForm(url, formValues)
	return getError(resp, err)
}

func getError(resp *http.Response, err error) error {
	if resp == nil || resp.Body == nil || err != nil {
		return err
	}
	defer resp.Body.Close()

	l := &loginReturn{}
	if err := json.NewDecoder(resp.Body).Decode(l); err != nil {
		return err
	}
	if l.ReturnCode != 0 {
		return fmt.Errorf("Error logging in: %s", l.ReturnMessage)
	}
	return nil
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
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

	resp, err := c.client.Get(url)
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

func checkResponse(returnCode int, message string) error {
	if returnCode != 0 {
		return fmt.Errorf("return code %d %s", returnCode, message)
	}
	return nil
}
