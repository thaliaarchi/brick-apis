package legobap

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LegoBAPClient struct {
	age     int
	country CountryCode
}

func NewClient(age int, country CountryCode) *LegoBAPClient {
	return &LegoBAPClient{age, country}
}

func (c *LegoBAPClient) GetPart(id string) (*ProductInformation, error) {
	url := "https://www.lego.com/en-US/service/rpservice/getitemordesign?itemordesignnumber=" + id + "&isSalesFlow=true"
	var part ProductInformation
	if err := c.doGet(url, &part); err != nil {
		return nil, err
	}
	return &part, nil
}

func (c *LegoBAPClient) GetSet(id string) (*ProductInformation, error) {
	url := "https://www.lego.com/en-US/service/rpservice/getproduct?productnumber=" + id + "&isSalesFlow=true"
	var set ProductInformation
	if err := c.doGet(url, &set); err != nil {
		return nil, err
	}
	return &set, nil
}

func (c *LegoBAPClient) doGet(url string, v interface{}) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	cookie := fmt.Sprintf(`csAgeAndCountry={"age":"%d","countrycode":"%s"}`, c.age, c.country)
	request.Header.Add("Cookie", cookie)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

type ProductInformation struct {
	Product                Product     `json:"Product"`
	Bricks                 []Brick     `json:"Bricks"`
	ImageBaseURL           string      `json:"ImageBaseUrl"`
	UnavailableInformation interface{} `json:"UnAvailableInformation"`
}

type Brick struct {
	ItemNo               int64              `json:"ItemNo"`
	ItemDescription      string             `json:"ItemDescr"`
	ColorLikeDescription string             `json:"ColourLikeDescr"`
	ColorDescription     string             `json:"ColourDescr"`
	MaingroupDescription string             `json:"MaingroupDescr"`
	Asset                string             `json:"Asset"`
	MaxQty               int64              `json:"MaxQty"`
	IP                   bool               `json:"Ip"`
	Price                float64            `json:"Price"`
	CurrencyID           string             `json:"CId"`
	SQty                 int64              `json:"SQty"`
	DesignID             int64              `json:"DesignId"`
	PriceStr             string             `json:"PriceStr"`
	PriceWithTaxStr      string             `json:"PriceWithTaxStr"`
	ItemUnavailable      bool               `json:"ItemUnavailable"`
	UnavailableLink      *UnavailableLink   `json:"UnavailableLink"`
	UnavailableReason    *UnavailableReason `json:"UnavailableReason"`
}

type UnavailableLink struct {
	URL   string `json:"Url"`
	Title string `json:"Title"`
}

type UnavailableReason struct {
	ReasonText        string        `json:"ReasonText"`
	LinkText          string        `json:"LinkText"`
	ID                string        `json:"ID"`
	Path              string        `json:"Path"`
	Key               string        `json:"Key"`
	HasVersion        bool          `json:"HasVersion"`
	RestrictedMarkets []interface{} `json:"RestrictedMarkets"`
}

type Product struct {
	ProductNo   string `json:"ProductNo"`
	ProductName string `json:"ProductName"`
	ItemNo      string `json:"ItemNo"`
	Asset       string `json:"Asset"`
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

/*
CId:
""
"USD"

ColourLikeDescr:
""
"Black"
"Blue"
"Grey"
"Purple"
"Red"
"White"
"Yellow"
*/

/*
MaingroupDescr:
"Bricks"
"Bricks, Special"
"Bricks, Special Circles And Angles"
*/

/*
ColorDescr:
"BLACK"
"BR.BLUE"
"BR.BLUEGREEN"
"BR.GREEN"
"BR.ORANGE"
"BR.RED"
"BR.YEL"
"BRICK-YEL"
"DK. BROWN"
"DK. ST. GREY"
"DK.GREEN"
"EARTH BLUE"
"FL. YELL-ORA"
"GOLD INK"
"LGH. PURPLE"
"M. LILAC"
"M. NOUGAT"
"MD.BLUE"
"MED. ST-GREY"
"MEDIUM AZUR"
"NEW DARK RED"
"OLIVE GREEN"
"RED. BROWN"
"SAND YELLOW"
"SILVER INK"
"SILVER MET."
"TR. BR. ORANGE"
"TR."
"TR.GREEN"
"TR.L.BLUE"
"TR.RED"
"W.GOLD"
"WHITE"
*/
