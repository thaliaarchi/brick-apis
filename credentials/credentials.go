package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Credentials struct {
	BrickLinkStore *BrickLinkStore `json:"bricklink_store"`
	BrickLinkUser  *BrickLinkUser  `json:"bricklink_user"`
	Brickset       *Brickset       `json:"brickset"`
	Lego           *LegoBAP        `json:"lego_bap"`
}

type BrickLinkStore struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	Token          string `json:"token"`
	TokenSecret    string `json:"token_secret"`
}

type BrickLinkUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Brickset struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

type LegoBAP struct {
	Age         string      `json:"age"`
	CountryCode CountryCode `json:"country_code"`
}

func Read(configFile string) (*Credentials, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var cred Credentials
	err = decoder.Decode(&cred)
	if err != nil {
		return nil, err
	}
	var errs []error
	if cred.BrickLinkUser.Username == "" {
		errs = append(errs, errors.New("BrickLink username must be set in credentials"))
	}
	if cred.BrickLinkUser.Password == "" {
		errs = append(errs, errors.New("BrickLink password must be set in credentials"))
	}
	if cred.BrickLinkStore.ConsumerKey == "" {
		errs = append(errs, errors.New("BrickLink consumer key must be set in credentials"))
	}
	if cred.BrickLinkStore.ConsumerSecret == "" {
		errs = append(errs, errors.New("BrickLink consumer secret must be set in credentials"))
	}
	if cred.BrickLinkStore.Token == "" {
		errs = append(errs, errors.New("BrickLink token must be set in credentials"))
	}
	if cred.BrickLinkStore.TokenSecret == "" {
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
	err = nil
	for i := range errs {
		if err == nil {
			err = errs[i]
		} else {
			err = fmt.Errorf("%s\n%s", err, errs[i])
		}
	}
	return &cred, err
}

// CountryCode is used by LEGO Bricks & Pieces.
type CountryCode string

// Only these countries can purchase through Bricks & Pieces.
const (
	CountryCodeAU CountryCode = "AU" // Australia
	CountryCodeAT CountryCode = "AT" // Austria
	CountryCodeBE CountryCode = "BE" // Belgium
	CountryCodeCA CountryCode = "CA" // Canada
	CountryCodeCZ CountryCode = "CZ" // Czech Republic
	CountryCodeDK CountryCode = "DK" // Denmark
	CountryCodeFI CountryCode = "FI" // Finland
	CountryCodeFR CountryCode = "FR" // France
	CountryCodeDE CountryCode = "DE" // Germany
	CountryCodeHU CountryCode = "HU" // Hungary
	CountryCodeIE CountryCode = "IE" // Ireland
	CountryCodeIT CountryCode = "IT" // Italy
	CountryCodeLU CountryCode = "LU" // Luxembourg
	CountryCodeNL CountryCode = "NL" // Netherlands
	CountryCodeNZ CountryCode = "NZ" // New Zealand
	CountryCodeNO CountryCode = "NO" // Norway
	CountryCodePL CountryCode = "PL" // Poland
	CountryCodePT CountryCode = "PT" // Portugal
	CountryCodeES CountryCode = "ES" // Spain
	CountryCodeSE CountryCode = "SE" // Sweden
	CountryCodeCH CountryCode = "CH" // Switzerland
	CountryCodeGB CountryCode = "GB" // United Kingdom
	CountryCodeUS CountryCode = "US" // United States
)
