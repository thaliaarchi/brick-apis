package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mrjones/oauth"
	"golang.org/x/net/publicsuffix"
)

var timeFormat = "2006/01/02 15:04:05"

type order struct {
	id   int
	time time.Time
}

func main() {
	username, password, consumerKey, consumerSecret, token, tokenSecret, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}

	userClient, err := createUserClient(username, password)
	if err != nil {
		log.Fatal(err)
	}
	exists, err := checkOrderExist(userClient, 9999999)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Order 9999999 exists: %t\n", exists)

	apiClient, err := createAPIClient(consumerKey, consumerSecret, token, tokenSecret)
	resp, err := apiClient.Get("https://api.bricklink.com/api/store/v1/orders/9999999")
	fmt.Println("Response:", resp.StatusCode, resp.Status)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Print(bodyString)
	}
}

func getCredentials() (string, string, string, string, string, string, error) {
	username := os.Getenv("BRICKLINK_USERNAME")
	if username == "" {
		return "", "", "", "", "", "", errors.New("BRICKLINK_USERNAME environment variable must be set")
	}
	password := os.Getenv("BRICKLINK_PASSWORD")
	if password == "" {
		return "", "", "", "", "", "", errors.New("BRICKLINK_PASSWORD environment variable must be set")
	}
	consumerKey := os.Getenv("BRICKLINK_CONSUMER_KEY")
	if consumerKey == "" {
		return "", "", "", "", "", "", errors.New("BRICKLINK_CONSUMER_KEY environment variable must be set")
	}
	consumerSecret := os.Getenv("BRICKLINK_CONSUMER_SECRET")
	if consumerSecret == "" {
		return "", "", "", "", "", "", errors.New("BRICKLINK_CONSUMER_SECRET environment variable must be set")
	}
	token := os.Getenv("BRICKLINK_TOKEN")
	if token == "" {
		return "", "", "", "", "", "", errors.New("BRICKLINK_TOKEN environment variable must be set")
	}
	tokenSecret := os.Getenv("BRICKLINK_TOKEN_SECRET")
	if tokenSecret == "" {
		return "", "", "", "", "", "", errors.New("BRICKLINK_TOKEN_SECRET environment variable must be set")
	}
	return username, password, consumerKey, consumerSecret, token, tokenSecret, nil
}

func createUserClient(username, password string) (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	_, err = client.PostForm("https://www.bricklink.com/ajax/renovate/loginandout.ajax", url.Values{
		"userid":          {username},
		"password":        {password},
		"keepme_loggedin": {"true"},
	})
	return client, err
}

func createAPIClient(consumerKey, consumerSecret, token, tokenSecret string) (*http.Client, error) {
	consumer := oauth.NewConsumer(consumerKey, consumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: token, Secret: tokenSecret}
	return consumer.MakeHttpClient(accessToken)
}

func checkOrderExist(client *http.Client, id int) (bool, error) {
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
