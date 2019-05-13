package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

const (
	renovateBase = "https://www.bricklink.com/ajax/renovate"
	cloneBase    = "https://www.bricklink.com/ajax/clone"
)

type BrickLinkUserClient struct {
	client      *http.Client
	credentials BrickLinkCredentials
}

func NewBrickLinkUserClient(cred *BrickLinkCredentials) (*BrickLinkUserClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	return &BrickLinkUserClient{&http.Client{Jar: jar}, *cred}, nil
}

func (c *BrickLinkUserClient) Login() error {
	_, err := c.client.PostForm(renovateBase+"/loginandout.ajax", url.Values{
		"userid":          {c.credentials.Username},
		"password":        {c.credentials.Password},
		"keepme_loggedin": {"true"},
	})
	return err
}

func (c *BrickLinkUserClient) GetWantedList(id int64) (*http.Response, error) {
	url := fmt.Sprintf(cloneBase+"/wanted/search2.ajax?wantedMoreID=%d", id)
	return c.client.Get(url)
}
