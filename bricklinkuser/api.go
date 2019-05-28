package bricklinkuser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/andrewarchi/brick-apis/credentials"
	"golang.org/x/net/publicsuffix"
)

const (
	renovateBase   = "https://www.bricklink.com/ajax/renovate"
	cloneBase      = "https://www.bricklink.com/ajax/clone"
	cloneStoreBase = "https://store.bricklink.com/ajax/clone"
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
	resp, err := c.client.PostForm(renovateBase+"/loginandout.ajax", url.Values{
		"userid":          {c.credentials.Username},
		"password":        {c.credentials.Password},
		"keepme_loggedin": {"true"},
	})
	return getError(resp, err)
}

func getError(resp *http.Response, err error) error {
	if resp == nil || resp.Body == nil || err != nil {
		return err
	}
	defer resp.Body.Close()

	l := &LoginReturn{}
	if err := json.NewDecoder(resp.Body).Decode(l); err != nil {
		return err
	}
	if l.ReturnCode != 0 {
		return fmt.Errorf("Error logging in: %s", l.ReturnMessage)
	}
	return nil
}

func (c *Client) GetWantedList(id int64) (*WantedListResults, error) {
	url := fmt.Sprintf(cloneBase+"/wanted/search2.ajax?wantedMoreID=%d", id)
	var wantedList wantedListResponse
	if err := c.doGet(url, &wantedList); err != nil {
		return nil, err
	}
	return &wantedList.Results, checkResponse(wantedList.ReturnCode, wantedList.ReturnMessage)
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func (c *Client) doGetAndSave(url string, v interface{}, filename string) error {
	file, err := os.Create("../data/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status %s", resp.Status)
	}
	defer resp.Body.Close()

	pr, pw := io.Pipe()
	tr := io.TeeReader(resp.Body, pw)

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
		decoder.DisallowUnknownFields()
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

func checkResponse(returnCode int, message string) error {
	if returnCode != 0 {
		return fmt.Errorf("return code %d %s", returnCode, message)
	}
	return nil
}

type wantedListResponse struct {
	Results        WantedListResults `json:"results"`
	ReturnCode     int               `json:"returnCode"`
	ReturnMessage  string            `json:"returnMessage"`
	ErrorTicket    int               `json:"errorTicket"`
	ProcessingTime int               `json:"procssingTime"`
}

type WantedListResults struct {
	ItemOptions    ItemOptions     `json:"itemOptions"`
	TotalResults   int             `json:"totalResults"`
	WantedLists    []WantedList    `json:"lists"`
	WantedItems    []WantedItem    `json:"wantedItems"`
	CategoryGroups []CategoryGroup `json:"categories"`
	ItemCount      int             `json:"totalCnt"`
	WantedListInfo WantedListInfo  `json:"wantedListInfo"`
	SearchMode     int             `json:"searchMode"`
	EmptySearch    int             `json:"emptySearch"`
}

type CategoryGroup struct {
	ItemType   ItemType       `json:"type"`
	Categories []CategoryInfo `json:"cats"`
	Total      int            `json:"total"`
}

type CategoryInfo struct {
	CategoryName string `json:"catName"`
	CategoryID   int    `json:"catID"`
	Count        int    `json:"cnt"`
	InvCount     int    `json:"invCnt"`
}

type ItemOptions struct {
	ShowStores bool `json:"showStores"`
}

type WantedList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WantedItem struct {
	WantedID          int             `json:"wantedID"`
	WantedListID      int             `json:"wantedMoreID"`
	WantedListName    string          `json:"wantedMoreName"`
	ItemNumber        string          `json:"itemNo"`
	ItemID            int             `json:"itemID"`
	ItemSeq           int             `json:"itemSeq"`
	ItemName          string          `json:"itemName"`
	ItemType          ItemType        `json:"itemType"`
	ItemBrand         int             `json:"itemBrand"`
	ImageURL          string          `json:"imgURL"`
	WantedQty         int             `json:"wantedQty"`
	WantedQtyFilled   int             `json:"wantedQtyFilled"`
	WantedCondition   WantedCondition `json:"wantedNew"`
	WantedNotify      WantedNotify    `json:"wantedNotify"`
	WantedRemark      string          `json:"wantedRemark"`
	WantedPrice       float64         `json:"wantedPrice"`
	FormatWantedPrice string          `json:"formatWantedPrice"`
	ColorID           int             `json:"colorID"`
	ColorName         string          `json:"colorName"`
	ColorHex          string          `json:"colorHex"`
}

type WantedListInfo struct {
	Name           string  `json:"name"`
	Description    string  `json:"desc"`
	ItemCount      int     `json:"num"`
	ID             int     `json:"id"`
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
