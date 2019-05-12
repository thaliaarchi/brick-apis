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
	Age            string
	CountryCode    string
}

func main() {
	os.Mkdir("data", 0777)

	cred, err := readCredentials("config.json")
	if err != nil {
		log.Fatal(err)
	}

	blUserClient, err := createBLUserClient(cred)
	if err != nil {
		log.Fatal(err)
	}
	blStoreClient, err := createBLStoreClient(cred)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := getBricksAndPiecesProduct(cred, "75192")
	writeResponse(resp, err, "set-75192.json")

	orders, err := getOrders(blStoreClient)
	if err != nil {
		fmt.Println(err)
	}
	if orders != nil {
		orders.printUnknownValues()
		for _, o := range orders.Orders {
			order, err := getOrderDetails(blStoreClient, o.OrderID)
			if err != nil {
				fmt.Println(err)
			}
			if order != nil {
				order.printUnknownValues()
			}
		}
	}

	resp, err = searchWantedList(blUserClient, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var search Search
	err = decodeAndWrite(resp.Body, &search, "wl-0.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(search)

	for _, list := range search.Results.WantedLists {
		if list.ID == 0 {
			continue
		}
		resp, err := searchWantedList(blUserClient, list.ID)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		var s Search
		err = decodeAndWrite(resp.Body, &s, fmt.Sprintf("wl-%d.json", list.ID))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
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
	if cred.Age == "" {
		return nil, errors.New("Age configuration variable must be set")
	}
	if cred.CountryCode == "" {
		return nil, errors.New("CountryCode configuration variable must be set")
	}
	return cred, nil
}

func printResponseCode(tag string, resp *http.Response) {
	fmt.Printf("%s: %d %s\n", tag, resp.StatusCode, resp.Status)
}

func printResponse(resp *http.Response, err error) {
	if err != nil {
		log.Fatal(err)
	}
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

func decodeAndWrite(r io.Reader, v interface{}, fileName string) error {
	file, err := os.Create("data/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	pr, pw := io.Pipe()
	tr := io.TeeReader(r, pw)

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
