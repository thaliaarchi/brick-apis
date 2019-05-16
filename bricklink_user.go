package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"

	"golang.org/x/net/publicsuffix"
)

const (
	renovateBase = "https://www.bricklink.com/ajax/renovate"
	cloneBase    = "https://www.bricklink.com/ajax/clone"
)

type BrickLinkUserClient struct {
	client      *http.Client
	credentials BrickLinkCredentials
}

func NewBrickLinkUserClient(cred *BrickLinkCredentials) (*BrickLinkUserClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	return &BrickLinkUserClient{&http.Client{Jar: jar}, *cred}, nil
}

func (c *BrickLinkUserClient) Login() error {
	_, err := c.client.PostForm(renovateBase+"/loginandout.ajax", url.Values{
		"userid":          {c.credentials.Username},
		"password":        {c.credentials.Password},
		"keepme_loggedin": {"true"},
	})
	return err
}

func (c *BrickLinkUserClient) GetWantedList(id int64) (*WantedListResults, error) {
	url := fmt.Sprintf(cloneBase+"/wanted/search2.ajax?wantedMoreID=%d", id)
	var response WantedListResponse
	if err := c.doRequest(url, fmt.Sprintf("Wanted List %d", id), fmt.Sprintf("wl-%d.json", id), &response); err != nil {
		return nil, err
	}
	return &response.Results, nil
}

func (c *BrickLinkUserClient) doRequest(url, tag, fileName string, v interface{}) error {
	resp, err := c.client.Get(url)
	printResponseCode(tag, resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := decodeAndWrite(resp.Body, v, fileName); err != nil {
		return err
	}
	f := reflect.ValueOf(v)
	if f != reflect.ValueOf(nil) {
		response := f.Interface().(*WantedListResponse)
		if response.ReturnMessage != "OK" {
			return fmt.Errorf("Response is not OK: %v", response)
		}
	}
	return nil
}

type WantedListResponse struct {
	Results        WantedListResults `json:"results"`
	ReturnCode     int64             `json:"returnCode"`
	ReturnMessage  string            `json:"returnMessage"`
	ErrorTicket    int64             `json:"errorTicket"`
	ProcessingTime int64             `json:"procssingTime"`
}

type WantedListResults struct {
	ItemOptions    ItemOptions     `json:"itemOptions"`
	TotalResults   int64           `json:"totalResults"`
	WantedLists    []WantedList    `json:"lists"`
	WantedItems    []WantedItem    `json:"wantedItems"`
	CategoryGroups []CategoryGroup `json:"categories"`
	ItemCount      int64           `json:"totalCnt"`
	WantedListInfo WantedListInfo  `json:"wantedListInfo"`
	SearchMode     int64           `json:"searchMode"`
	EmptySearch    int64           `json:"emptySearch"`
}

type CategoryGroup struct {
	ItemType   ItemType   `json:"type"`
	Categories []Category `json:"cats"`
	Total      int64      `json:"total"`
}

type Category struct {
	CatName  string `json:"catName"`
	CatID    int64  `json:"catID"`
	Count    int64  `json:"cnt"`
	InvCount int64  `json:"invCnt"`
}

type ItemOptions struct {
	ShowStores bool `json:"showStores"`
}

type WantedList struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WantedItem struct {
	WantedID          int64           `json:"wantedID"`
	WantedListID      int64           `json:"wantedMoreID"`
	WantedListName    string          `json:"wantedMoreName"`
	ItemNumber        string          `json:"itemNo"`
	ItemID            int64           `json:"itemID"`
	ItemSeq           int64           `json:"itemSeq"`
	ItemName          string          `json:"itemName"`
	ItemType          ItemType        `json:"itemType"`
	ItemBrand         int64           `json:"itemBrand"`
	ImageURL          string          `json:"imgURL"`
	WantedQty         int64           `json:"wantedQty"`
	WantedQtyFilled   int64           `json:"wantedQtyFilled"`
	WantedCondition   WantedCondition `json:"wantedNew"`
	WantedNotify      WantedNotify    `json:"wantedNotify"`
	WantedRemark      *string         `json:"wantedRemark"`
	WantedPrice       float64         `json:"wantedPrice"`
	FormatWantedPrice string          `json:"formatWantedPrice"`
	ColorID           int64           `json:"colorID"`
	ColorName         string          `json:"colorName"`
	ColorHex          string          `json:"colorHex"`
}

type WantedListInfo struct {
	Name           string  `json:"name"`
	Description    string  `json:"desc"`
	ItemCount      int64   `json:"num"`
	ID             int64   `json:"id"`
	CurrencySymbol string  `json:"curSymbol"`
	Completion     float64 `json:"progress"`
}

type ItemType string
type WantedCondition string
type WantedNotify string

const (
	ItemTypeSet         ItemType        = "S"
	ItemTypePart        ItemType        = "P"
	ItemTypeMinifig     ItemType        = "M"
	ItemTypeBook        ItemType        = "B"
	ItemTypeGear        ItemType        = "G"
	ItemTypeCatalog     ItemType        = "C"
	WantedConditionAny  WantedCondition = "X"
	WantedConditionNew  WantedCondition = "N"
	WantedConditionUsed WantedCondition = "U"
	WantedNotifyYes     WantedNotify    = "Y"
	WantedNotifyNo      WantedNotify    = "N"
)
