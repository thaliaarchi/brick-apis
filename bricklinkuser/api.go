package bricklinkuser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/google/go-querystring/query"
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
	return c.LoginAndOut(LoginAndOutOptions{
		Username:     username,
		Password:     password,
		StayLoggedIn: true,
	})
}

func (c *Client) Logout() error {
	return c.LoginAndOut(LoginAndOutOptions{DoLogout: true})
}

func (c *Client) LoginAndOut(options LoginAndOutOptions) error {
	url := fmt.Sprintf("https://%s/ajax/renovate/loginandout.ajax", getHost("www"))
	values, err := query.Values(options)
	if err != nil {
		return err
	}
	var r loginAndOutResult
	if err := c.doPost(url, &r, values); err != nil {
		return err
	}
	return checkResponse(r.ReturnCode, r.ReturnMessage)
}

type LoginAndOutOptions struct {
	Username     string `url:"userid,omitempty"`
	Password     string `url:"password,omitempty"`
	DoLogout     bool   `url:"do_logout,omitempty"`
	Override     bool   `url:"override,omitempty"`
	StayLoggedIn bool   `url:"keepme_loggedin,omitempty"`
	MID          string `url:"mid,omitempty"`
	PageID       string `url:"pageid,omitempty"` // Current page track ID
	Redirect     string `url:"redirect,omitempty"`
	ErrorURL     string `url:"error_url,omitempty"`
}

type loginAndOutResult struct {
	ReturnCode     int    `json:"returnCode"`
	ReturnMessage  string `json:"returnMessage"`
	ErrorTicket    int    `json:"errorTicket"`
	ProcessingTime int    `json:"procssingTime"`
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	return c.decode(resp, v)
}

func (c *Client) doPost(url string, v interface{}, formValues url.Values) error {
	resp, err := c.client.PostForm(url, formValues)
	if err != nil {
		return err
	}
	return c.decode(resp, v)
}

func (c *Client) decode(resp *http.Response, v interface{}) error {
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
