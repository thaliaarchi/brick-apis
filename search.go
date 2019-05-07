package main

import (
	"encoding/json"
	"io"
)

func DecodeSearch(r io.Reader) (Search, error) {
	var s Search
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&s)
	return s, err
}

type Search struct {
	Results        Results `json:"results"`
	ReturnCode     int64   `json:"returnCode"`
	ReturnMessage  string  `json:"returnMessage"`
	ErrorTicket    int64   `json:"errorTicket"`
	ProcessingTime int64   `json:"procssingTime"`
}

type Results struct {
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
	WantedID          int64        `json:"wantedID"`
	WantedMoreID      int64        `json:"wantedMoreID"`
	WantedMoreName    string       `json:"wantedMoreName"`
	ItemNo            string       `json:"itemNo"`
	ItemID            int64        `json:"itemID"`
	ItemSeq           int64        `json:"itemSeq"`
	ItemName          string       `json:"itemName"`
	ItemType          ItemType     `json:"itemType"`
	ItemBrand         int64        `json:"itemBrand"`
	ImageURL          string       `json:"imgURL"`
	WantedQty         int64        `json:"wantedQty"`
	WantedQtyFilled   int64        `json:"wantedQtyFilled"`
	WantedNew         WantedNew    `json:"wantedNew"`
	WantedNotify      WantedNotify `json:"wantedNotify"`
	WantedRemark      *string      `json:"wantedRemark"`
	WantedPrice       float64      `json:"wantedPrice"`
	FormatWantedPrice string       `json:"formatWantedPrice"`
	ColorID           int64        `json:"colorID"`
	ColorName         string       `json:"colorName"`
	ColorHex          string       `json:"colorHex"`
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
type WantedNew string
type WantedNotify string

const (
	TypeG         ItemType     = "G"
	TypeM         ItemType     = "M"
	TypeP         ItemType     = "P"
	TypeS         ItemType     = "S"
	WantedNewN    WantedNew    = "N"
	WantedNewX    WantedNew    = "X"
	WantedNotifyN WantedNotify = "N"
	WantedNotifyY WantedNotify = "Y"
)
