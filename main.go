package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// credentials stores user account and OAuth api credentials
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
	apiClient, err := createAPIClient(cred)
	if err != nil {
		log.Fatal(err)
	}

	printResponse(getOrderDetails(apiClient, 9999999))
	resp, err := searchWantedList(userClient, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	search, err := DecodeSearch(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	for _, list := range search.Results.WantedLists {
		resp, err := searchWantedList(userClient, list.ID)
		writeResponse(resp, err, fmt.Sprintf("wl-%d.json", list.ID))
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

func printResponse(resp *http.Response, err error) {
	fmt.Println("Response:", resp.StatusCode, resp.Status)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		str, err := responseToString(resp)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(str)
	}
}

func writeResponse(resp *http.Response, err error, fileName string) {
	fmt.Println("Response:", resp.StatusCode, resp.Status)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	io.Copy(file, resp.Body)
}

func responseToString(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), err
}
