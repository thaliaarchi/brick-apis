package bricklinkuser

import "fmt"

func (c *Client) GetWantedList(id int) (*WantedListResults, error) {
	url := fmt.Sprintf("https://%s/ajax/clone/wanted/search2.ajax?wantedMoreID=%d", getHost("www"), id)
	var w wantedListResponse
	if err := c.doGet(url, &w); err != nil {
		return nil, err
	}
	return &w.Results, checkResponse(w.ReturnCode, w.ReturnMessage)
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

type ItemOptions struct {
	ShowStores bool `json:"showStores"`
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
