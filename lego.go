package main

import (
	"fmt"
	"log"
	"net/http"
)

func getBricksAndPiecesPart(cred *credentials, id string) (*ProductInformation, error) {
	url := "https://www.lego.com/en-US/service/rpservice/getitemordesign?itemordesignnumber=" + id + "&isSalesFlow=true"
	return doLEGORequest(cred, url, fmt.Sprintf("Part %s", id), fmt.Sprintf("part-%s.json", id))
}

func getBricksAndPiecesSet(cred *credentials, id string) (*ProductInformation, error) {
	url := "https://www.lego.com/en-US/service/rpservice/getproduct?productnumber=" + id + "&isSalesFlow=true"
	return doLEGORequest(cred, url, fmt.Sprintf("Part %s", id), fmt.Sprintf("set-%s.json", id))
}

func doLEGORequest(cred *credentials, url, tag, fileName string) (*ProductInformation, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	cookie := fmt.Sprintf(`csAgeAndCountry={"age":"%s","countrycode":"%s"}`, cred.Age, cred.CountryCode)
	request.Header.Add("Cookie", cookie)
	resp, err := http.DefaultClient.Do(request)
	printResponseCode(tag, resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var product ProductInformation
	err = decodeAndWrite(resp.Body, &product, fileName)
	return &product, err
}

type ProductInformation struct {
	Product                Product     `json:"Product"`
	Bricks                 []Brick     `json:"Bricks"`
	ImageBaseURL           string      `json:"ImageBaseUrl"`
	UnAvailableInformation interface{} `json:"UnAvailableInformation"`
}

type Brick struct {
	ItemNo            int64              `json:"ItemNo"`
	ItemDescr         string             `json:"ItemDescr"`
	ColourLikeDescr   string             `json:"ColourLikeDescr"`
	ColourDescr       string             `json:"ColourDescr"`
	MaingroupDescr    string             `json:"MaingroupDescr"`
	Asset             string             `json:"Asset"`
	MaxQty            int64              `json:"MaxQty"`
	IP                bool               `json:"Ip"`
	Price             float64            `json:"Price"`
	CID               string             `json:"CId"`
	SQty              int64              `json:"SQty"`
	DesignID          int64              `json:"DesignId"`
	PriceStr          string             `json:"PriceStr"`
	PriceWithTaxStr   string             `json:"PriceWithTaxStr"`
	ItemUnavailable   bool               `json:"ItemUnavailable"`
	UnavailableLink   *UnavailableLink   `json:"UnavailableLink"`
	UnavailableReason *UnavailableReason `json:"UnavailableReason"`
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
