package bricklinkuser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/andrewarchi/brick-apis/credentials"
	"golang.org/x/net/publicsuffix"
)

const (
	renovateBase = "https://www.bricklink.com/ajax/renovate"
	cloneBase    = "https://www.bricklink.com/ajax/clone"
)

type Client struct {
	client      *http.Client
	credentials credentials.BrickLinkUser
}

func NewClient(cred *credentials.BrickLinkUser) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	return &Client{&http.Client{Jar: jar}, *cred}, nil
}

func (c *Client) Login() error {
	_, err := c.client.PostForm(renovateBase+"/loginandout.ajax", url.Values{
		"userid":          {c.credentials.Username},
		"password":        {c.credentials.Password},
		"keepme_loggedin": {"true"},
	})
	return err
}

func (c *Client) GetWantedList(id int64) (*WantedListResults, error) {
	url := fmt.Sprintf(cloneBase+"/wanted/search2.ajax?wantedMoreID=%d", id)
	var response WantedListResponse
	if err := c.doGet(url, &response); err != nil {
		return nil, err
	}
	return &response.Results, nil
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
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
