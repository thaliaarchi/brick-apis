package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"

	"github.com/mrjones/oauth"
	"golang.org/x/net/publicsuffix"
)

type credentials struct {
	Username       string
	Password       string
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

func main() {
	cred, err := readCredentials("config.json")
	if err != nil {
		log.Fatal(err)
	}

	userClient, err := createUserClient(cred)
	if err != nil {
		log.Fatal(err)
	}
	exists, err := checkOrderExist(userClient, 9999999)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Order 9999999 exists: %t\n", exists)

	apiClient, err := createAPIClient(cred)
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

func readCredentials(configFile string) (*credentials, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	cred := &credentials{}
	err = decoder.Decode(cred)
	if err != nil {
		return nil, err
	}
	if cred.Username == "" {
		return nil, errors.New("Username configuration variable must be set")
	}
	if cred.Password == "" {
		return nil, errors.New("Password configuration variable must be set")
	}
	if cred.ConsumerKey == "" {
		return nil, errors.New("ConsumerKey configuration variable must be set")
	}
	if cred.ConsumerSecret == "" {
		return nil, errors.New("ConsumerSecret configuration variable must be set")
	}
	if cred.Token == "" {
		return nil, errors.New("Token configuration variable must be set")
	}
	if cred.TokenSecret == "" {
		return nil, errors.New("TokenSecret configuration variable must be set")
	}
	return cred, nil
}

func createUserClient(cred *credentials) (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	_, err = client.PostForm("https://www.bricklink.com/ajax/renovate/loginandout.ajax", url.Values{
		"userid":          {cred.Username},
		"password":        {cred.Password},
		"keepme_loggedin": {"true"},
	})
	return client, err
}

func createAPIClient(cred *credentials) (*http.Client, error) {
	consumer := oauth.NewConsumer(cred.ConsumerKey, cred.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cred.Token, Secret: cred.TokenSecret}
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
