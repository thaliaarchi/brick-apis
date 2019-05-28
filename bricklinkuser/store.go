package bricklinkuser

import (
	"encoding/json"
	"io"
	"net/url"
)

// CartItemSimple is a simplified representation of an item used to add the item to a cart
type CartItemSimple struct {
	ID         int    `json:"invID"`
	Quantity   string `json:"invQty"`
	SellerID   int    `json:"sellerID"`
	SourceType int    `json:"sourceType"`
}

// AddToCartResponse is the reply given from the AddToCart method
type AddToCartResponse struct {
	Errors              int                `json:"errors"`
	ItemReturnStatus    []ItemReturnStatus `json:"itemReturnStatus"`
	Carts               []StoreCart        `json:"carts"`
	TotalStoreCartCount int                `json:"totStoreCartCnt"`
	CartItemErrorCode   int                `json:"cartItemErrorCode"`
	ReturnCode          int                `json:"returnCode"`
	ReturnMessage       string             `json:"returnMessage"`
	ErrorTicket         int                `json:"errorTicket"`
	ProcessingTime      int                `json:"procssingTime"`
}

// ItemReturnStatus provides information regarding itms to be returned
type ItemReturnStatus struct {
	InventoryID int    `json:"invID"`
	Code        string `json:"code"`
	Message     string `json:"msg"`
	SID         int    `json:"sid"`
}

// StoreCart groups together all the items a buyer would like to buy from a store
type StoreCart struct {
	SellerID    int       `json:"sellerID"`
	VATPct      int       `json:"vatPct"`
	SellerName  string    `json:"sellerName"`
	StoreName   string    `json:"storeName"`
	CountryID   string    `json:"countryID"`
	Feedback    int       `json:"feedback"`
	CurrentCart CartItems `json:"current_cart"`
}

// CartItems is the actual items in the StoreCart plus totals and other metadata
type CartItems struct {
	Items               []CartItemDetail `json:"items"`
	Superlots           []string         `json:"superlots"`
	TotalItems          int              `json:"totalItems"`
	TotalLots           int              `json:"totalLots"`
	TotalPrice          string           `json:"totalPrice"`
	TotalNativePrice    string           `json:"totalNativePrice"`
	TotalWarnings       int              `json:"totalWarnings"`
	TotalNativePriceRaw string           `json:"totalNativePriceRaw"`
	TotalWeightGrams    string           `json:"totalWeightGrams"`
	TotalWeightOunces   string           `json:"totalWeightOunces"`
	WeightUnknownLots   int              `json:"weightUnknownLots"`
	AverageLotPrice     string           `json:"aveLotPrice"`
}

// CartItemDetail provides detailed information about an item in a shopping cart
type CartItemDetail struct {
	ItemName                      string  `json:"itemName"`
	InventoryDescription          string  `json:"invDescription"`
	InventoryID                   int     `json:"invID"`
	InventoryQuantity             int     `json:"invQty"`
	BulkQuantity                  int     `json:"bulkQty"`
	SuperLotID                    int     `json:"superlotID"`
	SuperLotQuantity              int     `json:"superlotQty"`
	SalePercent                   int     `json:"salePercent"`
	ItemType                      string  `json:"itemType"`
	ItemBrand                     int     `json:"itemBrand"`
	InventoryCondition            string  `json:"invNew"`
	InventoryComplete             string  `json:"invComplete"`
	ColorID                       int     `json:"colorID"`
	ColorName                     string  `json:"colorName"`
	ItemNumber                    string  `json:"itemNo"`
	ItemSequence                  int     `json:"itemSeq"`
	ItemID                        int     `json:"itemID"`
	ItemStatus                    string  `json:"itemStatus"`
	SmallImage                    string  `json:"smallImg"`
	LargeImage                    string  `json:"largeImg"`
	NativePrice                   string  `json:"nativePrice"`
	SalePrice                     string  `json:"salePrice"`
	InventoryPrice                string  `json:"invPrice"`
	InventoryTierQuantity1        int     `json:"invTierQty1"`
	InventoryTierPrice1           string  `json:"invTierPrice1"`
	InventoryTierSalePrice1       string  `json:"invTierSalePrice1"`
	InventoryTierNativeSalePrice1 string  `json:"invTierNativeSalePrice1"`
	InventoryTierQuantity2        int     `json:"invTierQty2"`
	InventoryTierPrice2           string  `json:"invTierPrice2"`
	InventoryTierSalePrice2       string  `json:"invTierSalePrice2"`
	InventoryTierNativeSalePrice2 string  `json:"invTierNativeSalePrice2"`
	InventoryTierQuantity3        int     `json:"invTierQty3"`
	InventoryTierPrice3           string  `json:"invTierPrice3"`
	InventoryTierSalePrice3       string  `json:"invTierSalePrice3"`
	InventoryTierNativeSalePrice3 string  `json:"invTierNativeSalePrice3"`
	CartQuantity                  int     `json:"cartQty"`
	CartBindQuantity              int     `json:"cartBindQty"`
	InventoryDate                 string  `json:"invDate"`
	InventoryASCAvailable         bool    `json:"invASCAvailable"`
	InventoryAvailable            string  `json:"invAvailable"`
	Warnings                      []int   `json:"warnings"`
	TotalWeightOunces             float32 `json:"totalWeightOunces"`
	TotalWeightGrams              string  `json:"totalWeightGrams"`
	TotalPrice                    string  `json:"totalPrice"`
	TotalSalePrice                string  `json:"totalSalePrice"`
	TotalNativePrice              string  `json:"totalNativePrice"`
	TotalNativeSalePrice          string  `json:"totalNativeSalePrice"`
}

// AddToCart is used to add a list of items to a user's cart
func (c *Client) AddToCart(sid string, itemArray []CartItemSimple) (*AddToCartResponse, error) {
	q, err := getAddToCartQuery(sid, itemArray)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.PostForm(cloneBase+"/cart/add.ajax", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeCartReturn(resp.Body)
}

func getAddToCartQuery(sid string, itemArray []CartItemSimple) (url.Values, error) {
	values := url.Values{}
	data, err := json.Marshal(itemArray)
	if err != nil {
		return nil, err
	}
	values.Add("itemArray", string(data))
	values.Add("sid", sid)
	return values, nil
}

func decodeCartReturn(r io.ReadCloser) (*AddToCartResponse, error) {
	defer r.Close()
	retVal := &AddToCartResponse{}
	return retVal, json.NewDecoder(r).Decode(retVal)
}
