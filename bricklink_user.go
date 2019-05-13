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

func createBLUserClient(cred *BrickLinkCredentials) (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	_, err = client.PostForm(renovateBase+"/loginandout.ajax", url.Values{
		"userid":          {cred.Username},
		"password":        {cred.Password},
		"keepme_loggedin": {"true"},
	})
	return client, err
}

func searchWantedList(client *http.Client, id int64) (*http.Response, error) {
	url := fmt.Sprintf(cloneBase+"/wanted/search2.ajax?wantedMoreID=%d", id)
	return client.Get(url)
}
