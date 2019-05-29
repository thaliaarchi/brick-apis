package bricklinkuser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Converted from https://static.bricklink.com/_cache/jslegacy.2ec6747ecd9c6e44b6ac7e545e3f0457.js on 2019-05-25

// Converted from blUtil.refreshSession in jslegacy
func refreshSession() {
}

// Converted from blUtil.getRand2Byte in jslegacy
func getRand2Byte() string {
	return strconv.FormatInt(int64((rand.Float64()+1.0)*0x10000), 16)[1:]
}

// Converted from blUtil.createMID in jslegacy
func createMID() string {
	ms := time.Now().UnixNano() / 1000000
	if ms < 0 {
		ms = -ms
	}
	msStr := (strconv.FormatInt(ms, 16) + "0000000000000000")[0:16]
	mid := msStr + "-" + getRand2Byte() + getRand2Byte() + getRand2Byte() + getRand2Byte()
	return strings.ToLower(mid)
}

// getBLHost returns the host URL for a given type.
// Converted from blUtil.getBLHost in jslegacy.
func getBLHost(hostType string) string {
	// JS would override with values from blo_host variable, but I cannot find those values
	host := "www.bricklink.com"
	switch hostType {
	case "www", "alpha":
		host = "www.bricklink.com"
	case "img":
		host = "img.bricklink.com"
	case "static":
		host = "static.bricklink.com"
	case "store":
		host = "store.bricklink.com"
	}
	return host
}

// Converted from blURL.getCatalogItemPageURL in jslegacy.
func getCatalogItemPageURL(itemID int) string {
	return fmt.Sprintf("//%s/v2/catalog/catalogitem.page?id=%d", getBLHost("www"), itemID)
}

// Converted from blURL.getCatalogItemPageURLWithColor in jslegacy.
func getCatalogItemPageURLWithColor(itemID, colorID int) string {
	return fmt.Sprintf("//%s/v2/catalog/catalogitem.page?id=%d&idColor=%d", getBLHost("www"), itemID, colorID)
}

// Converted from blURL.getCatalogItemPageURLByItemNo in jslegacy.
func getCatalogItemPageURLByItemNo(itemType rune, itemNo string, itemSeq int) string {
	seq := ""
	if itemType == 'S' || itemType == 'I' || itemType == 'O' {
		seq = "-" + strconv.Itoa(itemSeq)
	}
	return fmt.Sprintf("//%s/v2/catalog/catalogitem.page?%c=%s%s", getBLHost("www"), itemType, itemNo, seq)
}

// Converted from blURL.getCatalogItemsForSalePageURL in jslegacy.
func getCatalogItemsForSalePageURL(itemID, colorID int) string {
	return fmt.Sprintf("//%s/v2/catalog/catalogitem.page?id=%d&idColor=%d#T=S", getBLHost("www"), itemID, colorID)
}

// Converted from blURL.getNewsPageURL in jslegacy.
func getNewsPageURL(msgID int) string {
	return fmt.Sprintf("//%s/v2/community/newsview.page?msgid=%d", getBLHost("www"), msgID)
}

// Converted from blURL.getStoreURL and equivalent to blURL.getNewStoreURL in jslegacy.
func getStoreURL(sellerUsername string) string {
	return fmt.Sprintf("//%s/%s", getBLHost("store"), sellerUsername)
}

// Converted from blURL.getStoreURLByID in jslegacy.
func getStoreURLByID(sellerUserID int) string {
	return fmt.Sprintf("//%s/store/home.page?sid=%d", getBLHost("www"), sellerUserID)
}

// Converted from blURL.getStoreInvURL in jslegacy.
func getStoreInvURL(sellerUsername string, invID int) string {
	return fmt.Sprintf("//%s/%s?itemID=%d", getBLHost("store"), sellerUsername, invID)
}

// Converted from blURL.getStoreInvURLByID in jslegacy.
func getStoreInvURLByID(sellerUserID, invID int) string {
	return fmt.Sprintf("//%s/store/home.page?sid=%d&itemID=%d", getBLHost("www"), sellerUserID, invID)
}

type storeOptions struct {
	Query            string  `json:"q,omitempty"`      // Search
	Sort             int     `json:"sort,omitempty"`   // 1: Item Name, Color; 2: Item Number, Color; 3: Condition, Item Name; 4: Color, Item Name; 6: Price; 7: Quantity; 8: Sale Amount; 9: Date Added
	Descending       int     `json:"desc,omitempty"`   // 0: false, 1: true
	PageSize         int     `json:"pgSize,omitempty"` // 10, 25, 50, 100
	Page             int     `json:"pg,omitempty"`
	InvID            string  `json:"invID,omitempty"`
	ItemID           string  `json:"itemID,omitempty"`
	ItemType         string  `json:"itemType,omitempty"`
	ItemBrandFilter  string  `json:"itemBrandFilter,omitempty"`
	ItemTypeFilter   string  `json:"itemTypeFilter,omitempty"` // "": All Item Types, "S": Sets, "P": Parts, "M": Minifigs, "B": Books, "G": Gear, "C": Catalogs, "I": Instructions, "O": Original Boxes, "U": Unsorted Lots
	CategoryID       string  `json:"catID,omitempty"`
	CategoryIDFilter string  `json:"catIDFilter,omitempty"`
	ItemYear         int     `json:"itemYear,omitempty,string"`
	ColorID          string  `json:"colorID,omitempty"`
	ColorIDFilter    string  `json:"colorIDFilter,omitempty"`
	WantedListIDs    string  `json:"wantedMoreArrayID,omitempty"` // comma separated list
	ReservedUserID   int     `json:"resUserID,omitempty"`
	QuantityMin      int     `json:"Qmin,omitempty,string"`
	QuantityMax      int     `json:"Qmax,omitempty,string"`
	PriceMin         float64 `json:"Pmin,omitempty,string"`      // Price in store currency
	PriceMax         float64 `json:"Pmax,omitempty,string"`      // Price in store currency
	OnSale           int     `json:"bOnSale,omitempty"`          // 1: Show Items on Sale
	OnlyCustomItems  int     `json:"bOnlyCustomItems,omitempty"` // 0: show, 1: hide, 2: only
	ExcludeSuperLot  int     `json:"bExcludeSuperLot,omitempty"` // 0: show, 1: hide, 2: only
	ExcludeTiered    int     `json:"bExcludeTiered,omitempty"`   // 0: show, 1: hide, 2: only
	ExcludeBulk      int     `json:"bExcludeBulk,omitempty"`     // 0: show, 1: hide, 2: only
	InvNew           string  `json:"invNew,omitempty"`           // "": all, "N": new, "U": used
	ItemStatus       string  `json:"itemStatus,omitempty"`
	BindType         string  `json:"bindType,omitempty"`
	BindID           string  `json:"bindID,omitempty"`
	OnWantedList     int     `json:"bOnWantedList,omitempty"` // 0: false, 1: true
	ShowHomeItems    int     `json:"showHomeItems,omitempty"` // 0: false, 1: true (Featured)
	ShowNewest       int     `json:"showNewest,omitempty"`    // 0: false, 1: true
	HideHaveMore     int     `json:"bHideHaveMore,omitempty"` // 1: Hide Items if Have Qty is ≥ Want Qty
}

// Converted from blURL.getStoreWLURL in jslegacy.
func getStoreWLURL(sellerUsername string, wantedListIDs []int) string {
	options := storeOptions{
		OnWantedList:  1,
		WantedListIDs: sliceJoin(wantedListIDs, ","),
	}
	bytes, err := json.Marshal(&options)
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("//%s/%s#/shop?o=%s", getBLHost("store"), sellerUsername, string(bytes))
}

func sliceJoin(slice []int, delim string) string {
	var buffer bytes.Buffer
	for i := range slice {
		if i != 0 {
			buffer.WriteString(delim)
		}
		buffer.WriteString(strconv.Itoa(slice[i]))
	}
	return buffer.String()
}

// Converted from blURL.getStoreCartURL in jslegacy.
func getStoreCartURL(sellerUsername string) string {
	return fmt.Sprintf("//%s/%s#/cart", getBLHost("store"), sellerUsername)
}

// Converted from blURL.getStoreCheckoutURL in jslegacy.
func getStoreCheckoutURL(sellerUsername string) string {
	return fmt.Sprintf("//%s/%s#/checkout", getBLHost("store"), sellerUsername)
}

// Converted from blURL.getStoreCartURLByID in jslegacy.
func getStoreCartURLByID(sellerUserID int) string {
	return fmt.Sprintf("//%s/store/home.page?sid=%d#/cart", getBLHost("www"), sellerUserID)
}

// Acceptable values for flagSize are 'S', 'M', and 'L'
// Generalized from blURL.getCountryFlagSmallURL and blURL.getCountryFlagMediumURL in jslegacy.
func getCountryFlagURL(countryID string, flagSize rune) string {
	return fmt.Sprintf("//%s/Images/Flags%c/%s.gif", getBLHost("img"), flagSize, countryID)
}

// Converted from blURL.getStoreFeedbackURL in jslegacy.
func getStoreFeedbackURL(sellerUsername string) string {
	return fmt.Sprintf("//%s/store/home.page?p=%s#/feedback", getBLHost("www"), sellerUsername)
}

// Converted from blURL.getLoginURL in jslegacy.
func getLoginURL(loginTo string) string {
	return fmt.Sprintf("https://%s/v2/login.page?logInTo=%s", getBLHost("www"), url.QueryEscape(loginTo))
}

// Converted from blURL.getDefaultStoreLogoURL in jslegacy.
func getDefaultStoreLogoURL() string {
	return fmt.Sprintf("//%s/clone/img/store-default-image.png", getBLHost("static"))
}

// Converted from blURL.getFeedbackIconUrl in jslegacy.
func getFeedbackIconURL(score int) string {
	index := "000"
	if score < 10 {
		index = "000"
	} else if score < 50 {
		index = "001"
	} else if score < 100 {
		index = "002"
	} else if score < 500 {
		index = "003"
	} else if score < 1000 {
		index = "004"
	} else if score < 2500 {
		index = "005"
	} else if score < 5000 {
		index = "006"
	} else if score < 10000 {
		index = "007"
	} else if score < 25000 {
		index = "008"
	} else if score < 50000 {
		index = "009"
	} else {
		index = "010"
	}
	return fmt.Sprintf("//static.bricklink.com/clone/img/feedback_%s.png", index)
}

// Converted from blc_GlobalCart.retrieveCartInfo
func (c *Client) GetGlobalCart() (*CartInfo, error) {
	url := fmt.Sprintf("https://%s/ajax/renovate/getglobalcart.ajax", getBLHost("www"))
	var cartInfo CartInfo
	if err := c.doGet(url, &cartInfo); err != nil {
		return nil, err
	}
	return &cartInfo, checkResponse(cartInfo.ReturnCode, cartInfo.ReturnMessage)
}

type CartInfo struct {
	List            []StoreList `json:"list"`
	TotalStoreCount int         `json:"total_store_cnt"`
	TotalLotCount   int         `json:"total_lot_cnt"`
	ReturnCode      int         `json:"returnCode"`
	ReturnMessage   string      `json:"returnMessage"`
	ErrorTicket     int         `json:"errorTicket"`
	ProcessingTime  int         `json:"procssingTime"`
}

type StoreList struct {
	SellerID        string  `json:"sellerid"`
	StoreName       string  `json:"store_name"`
	Username        string  `json:"username"`
	CountryID       string  `json:"countryid"`
	FeedbackScore   int64   `json:"feedback_score"`
	InstantCheckout bool    `json:"instantCheckout"`
	LotCount        int64   `json:"lotcnt"`
	DispPrice       float64 `json:"fDispPrice"`
	TotalPrice      string  `json:"strTotPrice"`
	Key             string  `json:"key"`
}

// Converted from blc_GlobalCart.getCheckoutInfo
func (c *Client) GetGlobalCartCheckoutInfo(sellerUserID int, key string) (*CheckoutInfo, error) {
	// Single store: (also conditions.estShippingAndHandling)
	// { action: 'conditions', sid: store.sellerid, key: store.key, checkPaypal: 0 }
	url := fmt.Sprintf("https://%s/ajax/clone/store/preparecheckout.ajax?action=conditions&sid=%d&key=%s&checkPaypal=0", getBLHost("www"), sellerUserID, key)
	var checkoutInfo CheckoutInfo
	if err := c.doGetAndSave(url, checkoutInfo, "checkout-info.json"); err != nil {
		return nil, err
	}
	return &checkoutInfo, nil
}

type CheckoutInfo struct {
	Conditions     Conditions `json:"conditions"`
	ReturnCode     int        `json:"returnCode"`
	ReturnMessage  string     `json:"returnMessage"`
	ErrorTicket    int        `json:"errorTicket"`
	ProcessingTime int        `json:"procssingTime"`
}

type A struct {
	Sellers        []SellerConditions `json:"sellers"`
	ReturnCode     int64              `json:"returnCode"`
	ReturnMessage  string             `json:"returnMessage"`
	ErrorTicket    int64              `json:"errorTicket"`
	ProcessingTime int64              `json:"procssingTime"`
}

type SellerConditions struct {
	SellerID   int64      `json:"sid"`
	Conditions Conditions `json:"conditions"`
}

type Conditions struct {
	Error                              ConditionsError      `json:"error"`
	TargetShippingMethod               TargetShippingMethod `json:"targetShippingMethod"`
	EstimatedShippingAndHandlingNative string               `json:"estShippingAndHandlingNative"`
	CostType                           int                  `json:"costType"`
	HasShippingCost                    bool                 `json:"hasShippingCost"`
	ShippingLocation                   string               `json:"shippingLocation"`
	Domestic                           bool                 `json:"domestic"`
	OnlyMethod                         bool                 `json:"onlyMethod"`
	OrderRestriction                   OrderRestriction     `json:"orderRestriction"`
	ShowSalesTax                       bool                 `json:"showSalesTax"`
	SalesTaxFinal                      bool                 `json:"salesTaxFinal"`
	BatchID                            int                  `json:"batchID"`
	BatchNum                           int                  `json:"batchNum"`
	BatchTotalNativePrice              string               `json:"batchTotalNativePrice"`
	OrderItemTotalNativePrice          string               `json:"orderItemTotalNativePrice"`
	AvgLotNativePrice                  string               `json:"avgLotNativePrice"`
	OrderTotalNativePrice              string               `json:"orderTotalNativePrice"`
	BatchTotalPrice                    string               `json:"batchTotalPrice"`
	OrderTotalPrice                    string               `json:"orderTotalPrice"`
	OrderItemTotalPrice                string               `json:"orderItemTotalPrice"`
	AvgLotPrice                        string               `json:"avgLotPrice"`
	LoggedIn                           bool                 `json:"loggedIn"`
	EmailConfirmed                     bool                 `json:"emailConfirmed"`
	SalesTaxNative                     string               `json:"salesTaxNative"`
	SalesTax                           string               `json:"salesTax"`
}

type ConditionsError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type OrderRestriction struct {
	MinBuyMet            bool   `json:"minBuyMet"`
	AvgBuyMet            bool   `json:"avgBuyMet"`
	MinBuyShortageNative string `json:"minBuyShortageNative"`
	AvgBuyShortageNative string `json:"avgBuyShortageNative"`
	MinBuyShortage       string `json:"minBuyShortage"`
	AvgBuyShortage       string `json:"avgBuyShortage"`
}

type TargetShippingMethod struct {
	ID                  int64                `json:"id"`
	Name                string               `json:"name"`
	UnitType            int64                `json:"unitType"`
	APIMethod           bool                 `json:"apiMethod"`
	PackageRestrictions []PackageRestriction `json:"packageRestrictions"`
}

type PackageRestriction struct {
	Type int64  `json:"type"`
	Arg1 string `json:"arg1"`
	Arg2 string `json:"arg2"`
	Arg3 string `json:"arg3"`
}

// blc_CatalogItem.getItemTypeName
func getItemTypeName(itemType rune, isPlural bool) string {
	switch itemType {
	case 'P', 'p':
		if isPlural {
			return "Parts"
		}
		return "Part"
	case 'S', 's':
		if isPlural {
			return "Sets"
		}
		return "Set"
	case 'M', 'm':
		if isPlural {
			return "Minifigs"
		}
		return "Minifig"
	case 'G', 'g':
		return "Gear"
	case 'I', 'i':
		if isPlural {
			return "Instructions"
		}
		return "Instruction"
	case 'O', 'o':
		if isPlural {
			return "Original Boxes"
		}
		return "Original Box"
	case 'C', 'c':
		if isPlural {
			return "Catalogs"
		}
		return "Catalog"
	case 'B', 'b':
		if isPlural {
			return "Books"
		}
		return "Book"
	case 'U', 'u':
		if isPlural {
			return "Custom Items"
		}
		return "Custom Item"
	default:
		return ""
	}
}
