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
	"strconv"
)

// credentials stores user account and OAuth api credentials
type Credentials struct {
	BrickLink BrickLinkCredentials `json:"bricklink"`
	Brickset  BricksetCredentials  `json:"brickset"`
	Lego      LegoCredentials      `json:"lego"`
}
type BrickLinkCredentials struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	Token          string `json:"token"`
	TokenSecret    string `json:"token_secret"`
}
type BricksetCredentials struct {
	Key string `json:"key"`
}
type LegoCredentials struct {
	Age         string `json:"age"`
	CountryCode string `json:"country_code"`
}

func main() {
	os.Mkdir("data", 0755)

	cred, errs := readCredentials("credentials.json")
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	blUser, err := NewBrickLinkUserClient(&cred.BrickLink)
	if err != nil {
		log.Fatal(err)
	}
	blUser.Login()
	blStore, err := NewBrickLinkStoreClient(&cred.BrickLink)
	if err != nil {
		log.Fatal(err)
	}
	lego := NewLegoClient(&cred.Lego)

	_, err = blStore.GetColorList()
	if err != nil {
		fmt.Println(err)
	}
	// for _, c := range colors {
	// 	_, err := getColor(blStoreClient, c.ColorID)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	orders, err := blStore.GetOrderList()
	if err != nil {
		fmt.Println(err)
	}
	for _, o := range orders {
		order, err := blStore.GetOrder(o.OrderID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(order)
		if order != nil {
			order.printUnknownValues()
		}
		o.printUnknownValues()
	}

	resp, err := blUser.GetWantedList(0)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var search Search
	err = decodeAndWrite(resp.Body, &search, "wl-0.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, list := range search.Results.WantedLists {
		if list.ID == 0 {
			continue
		}
		resp, err := blUser.GetWantedList(list.ID)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		var s Search
		err = decodeAndWrite(resp.Body, &s, fmt.Sprintf("wl-%d.json", list.ID))
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = lego.GetBricksAndPiecesPart("3024")
	if err != nil {
		fmt.Println(err)
	}
	_, err = lego.GetBricksAndPiecesSet("75192")
	if err != nil {
		fmt.Println(err)
	}
}

func readCredentials(configFile string) (*Credentials, []error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, []error{err}
	}
	decoder := json.NewDecoder(file)
	var cred Credentials
	err = decoder.Decode(&cred)
	if err != nil {
		return nil, []error{err}
	}
	var errs []error
	if cred.BrickLink.Username == "" {
		errs = append(errs, errors.New("BrickLink username must be set in credentials"))
	}
	if cred.BrickLink.Password == "" {
		errs = append(errs, errors.New("BrickLink password must be set in credentials"))
	}
	if cred.BrickLink.ConsumerKey == "" {
		errs = append(errs, errors.New("BrickLink consumer key must be set in credentials"))
	}
	if cred.BrickLink.ConsumerSecret == "" {
		errs = append(errs, errors.New("BrickLink consumer secret must be set in credentials"))
	}
	if cred.BrickLink.Token == "" {
		errs = append(errs, errors.New("BrickLink token must be set in credentials"))
	}
	if cred.BrickLink.TokenSecret == "" {
		errs = append(errs, errors.New("BrickLink token secret must be set in credentials"))
	}
	if cred.Brickset.Key == "" {
		errs = append(errs, errors.New("Brickset key must be set in credentials"))
	}
	if cred.Lego.Age == "" {
		errs = append(errs, errors.New("Age must be set in credentials"))
	}
	if cred.Lego.CountryCode == "" {
		errs = append(errs, errors.New("Country code must be set in credentials"))
	}
	age, err := strconv.Atoi(cred.Lego.Age)
	if err != nil {
		errs = append(errs, err)
	}
	if age < 18 {
		errs = append(errs, errors.New("Age must be at least 18 for Bricks & Pieces"))
	}
	switch cred.Lego.CountryCode {
	case "AU", "AT", "BE", "CA", "CZ", "DK", "FI", "FR", "DE", "HU", "IE",
		"IT", "LU", "NL", "NZ", "NO", "PL", "PT", "ES", "SE", "CH", "GB", "US":
	default:
		errs = append(errs, errors.New("Country is not supported for Bricks & Pieces"))
	}
	return &cred, errs
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
