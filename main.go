package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/publicsuffix"
)

var timeFormat = "2006/01/02 15:04:05"

type order struct {
	id   int
	time time.Time
}

func main() {
	username, password, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}
	createClient(username, password)
}

func getCredentials() (string, string, error) {
	username := os.Getenv("BRICKLINK_USERNAME")
	if username == "" {
		return "", "", errors.New("BRICKLINK_USERNAME environment variable must be set")
	}
	password := os.Getenv("BRICKLINK_PASSWORD")
	if password == "" {
		return "", "", errors.New("BRICKLINK_PASSWORD environment variable must be set")
	}
	return username, password, nil
}

func createClient(username, password string) http.Client {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}
	fmt.Println("Logging in")
	_, err = client.PostForm("https://www.bricklink.com/ajax/renovate/loginandout.ajax", url.Values{
		"userid":          {username},
		"password":        {password},
		"keepme_loggedin": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in")
	return client
}

func checkOrderExist(client http.Client, id int) (bool, error) {
	url := "https://www.bricklink.com/orderDetail.asp?ID=" + strconv.Itoa(id)
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}

	finalURL := resp.Request.URL.String()
	if finalURL == url || finalURL == "https://www.bricklink.com/oops.asp?err=403" {
		return true, nil
	} else if finalURL == "https://www.bricklink.com/notFound.asp?nf=order&mFolder=o&mSub=o" {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected URL: %v", finalURL)
	}
}
