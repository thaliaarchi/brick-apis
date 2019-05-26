package bricklinkuser

import "fmt"

// https://www.bricklink.com/ajax/clone/search/searchproduct.ajax?q=75159&st=0&cond=&brand=1000&type=&cat=&yf=0&yt=0&loc=&reg=0&ca=0&ss=&pmt=&nmp=0&color=-1&min=0&max=0&minqty=0&nosuperlot=1&incomplete=0&showempty=1&rpp=25&pi=1&ci=0

func (c *Client) SearchProduct() (*SearchProduct, error) {
	url := fmt.Sprintf("%s/search/searchproduct.ajax", cloneBase) // params TBD
	var searchProduct SearchProduct
	if err := c.doGet(url, &searchProduct); err != nil {
		return nil, err
	}
	return &searchProduct, checkResponse(searchProduct.ReturnCode, searchProduct.ReturnMessage)
}

type SearchProduct struct {
	TotalCount     int           `json:"total_count"`
	ColorID        int           `json:"idColor"`
	ResultsPerPage int           `json:"rpp"`
	PageIndex      int           `json:"pi"`
	List           []ProductList `json:"list"`
	ReturnCode     int           `json:"returnCode"`
	ReturnMessage  string        `json:"returnMessage"`
	ErrorTicket    int           `json:"errorTicket"`
	ProcessingTime int           `json:"procssingTime"`
}

type ProductList struct {
	InvID                  int          `json:"idInv"`
	Description            string       `json:"strDesc"`
	NewOrUsed              NewOrUsed    `json:"codeNew"`
	Completeness           Completeness `json:"codeComplete"`
	Quantity               int          `json:"n4Qty"`
	ColorID                int          `json:"idColor"`
	ColorIDDefault         int          `json:"idColorDefault"`
	Color                  string       `json:"strColor"`
	ImageURL               string       `json:"strInvImgUrl"`           // URL of image provided by seller
	ImageID                int          `json:"idInvImg"`               // ID of Image provided by seller
	ImageType              ImageType    `json:"typeInvImg"`             // Type of Image provided by seller
	ImageTypeDefault       ImageType    `json:"typeImgDefault"`         // Default catalog image type
	HasExtendedDescription int          `json:"hasExtendedDescription"` // Has an extended description (0: false, 1: true)
	InstantCheckout        bool         `json:"instantCheckout"`
	SalePrice              string       `json:"mInvSalePrice"`     // Sale price in currency of store
	SalePriceDisplay       string       `json:"mDisplaySalePrice"` // Sale price in display currency of user
	SalePercent            int          `json:"nSalePct"`          // Percent discounted
	Tier1Quantity          int          `json:"nTier1Qty"`
	Tier2Quantity          int          `json:"nTier2Qty"`
	Tier3Quantity          int          `json:"nTier3Qty"`
	Tier1Price             string       `json:"nTier1InvPrice"`     // Price in currency of store
	Tier2Price             string       `json:"nTier2InvPrice"`     // Price in currency of store
	Tier3Price             string       `json:"nTier3InvPrice"`     // Price in currency of store
	Tier1PriceDisplay      string       `json:"nTier1DisplayPrice"` // Price in display currency of user
	Tier2PriceDisplay      string       `json:"nTier2DisplayPrice"` // Price in display currency of user
	Tier3PriceDisplay      string       `json:"nTier3DisplayPrice"` // Price in display currency of user
	Category               string       `json:"strCategory"`
	StoreName              string       `json:"strStorename"`
	StoreCurrencyID        int          `json:"idCurrencyStore"`
	MinBuy                 string       `json:"mMinBuy"`
	SellerUsername         string       `json:"strSellerUsername"`
	SellerFeedbackScore    int          `json:"n4SellerFeedbackScore"`
	SellerCountryName      string       `json:"strSellerCountryName"`
	SellerCountryCode      string       `json:"strSellerCountryCode"`
}

type Completeness string

const (
	B Completeness = "B"
	C Completeness = "C"
	S Completeness = "S"
	X Completeness = "X"
)

type NewOrUsed string

const (
	N NewOrUsed = "N"
	U NewOrUsed = "U"
)

type ImageType string

const (
	ImageTypeEmpty ImageType = ""
	ImageTypeJ     ImageType = "J"
)

// GetStoreItem retrieves details for an item sold in a store.
// This API is called when clicking on an item in a store to show the details modal.
// The type information is currently incomplete and no examples have been found using the URL parameter wantedMoreArrayID.
func (c *Client) GetStoreItem(invID, storeID, wantedListArrayID string) (*StoreItem, error) {
	url := fmt.Sprintf("%s/store/item.ajax?invID=%s&sid=%s&wantedMoreArrayID=%s", cloneStoreBase, invID, storeID, wantedListArrayID)
	var storeItem StoreItem
	if err := c.doGet(url, &storeItem); err != nil {
		return nil, err
	}
	return &storeItem, checkResponse(storeItem.ReturnCode, storeItem.ReturnMessage)
}

type StoreItem struct {
	InvID               int                     `json:"invID"`
	Description         string                  `json:"description"`
	InvDescription      string                  `json:"invDescription"` // Found in React, not confirmed to be in json
	ExtendedDescription string                  `json:"extDescription"`
	InvQuantity         int                     `json:"invQty"`
	InvBulk             int                     `json:"invBulk"`
	InvSale             int                     `json:"invSale"`
	BulkQuantity        int                     `json:"bulkQty"` // Found in React, not confirmed to be in json
	ItemType            string                  `json:"itemType"`
	InvNew              string                  `json:"invNew"` // "N": New, "U": Used
	ColorID             int                     `json:"colorID"`
	ColorName           string                  `json:"colorName"`
	ColorHex            string                  `json:"colorHex"`
	InvDate             string                  `json:"invDate"`
	ItemNo              string                  `json:"itemNo"`
	ItemSeq             int                     `json:"itemSeq"`
	ItemID              int                     `json:"itemID"`
	ItemName            string                  `json:"itemName"`
	ItemStatus          string                  `json:"itemStatus"`
	ItemBrand           int                     `json:"itemBrand"`
	InvComplete         string                  `json:"invComplete"`
	NativePrice         string                  `json:"nativePrice"`
	SalePrice           string                  `json:"salePrice"`
	SalePercent         int                     `json:"salePercent"` // Found in React, not confirmed to be in json
	InvPrice            string                  `json:"invPrice"`
	InvTier1Quantity    int                     `json:"invTierQty1"`
	InvTier2Quantity    int                     `json:"invTierQty2"`
	InvTier3Quantity    int                     `json:"invTierQty3"`
	InvTier1Price       string                  `json:"invTierPrice1"`
	InvTier2Price       string                  `json:"invTierPrice2"`
	InvTier3Price       string                  `json:"invTierPrice3"`
	InvTier1NativePrice string                  `json:"invTierNativePrice1"`
	InvTier2NativePrice string                  `json:"invTierNativePrice2"`
	InvTier3NativePrice string                  `json:"invTierNativePrice3"`
	CartQuantity        int                     `json:"cartQty"`
	CartBindQuantity    int                     `json:"cartBindQty"`
	InvBindID           int                     `json:"invBindID"`
	InvBindQuantity     int                     `json:"invBindQty"`
	InvImageID          int                     `json:"invImgID"`
	InvImageType        string                  `json:"invImgType"`
	InvURL              string                  `json:"invURL"`
	ImageURL            string                  `json:"imgURL"`
	NewItemImage        string                  `json:"newItemImg"` // Found in React, not confirmed to be in json
	LargeImages         []LargeImage            `json:"largeImgs"`
	CategoryInfo        []StoreItemCategoryInfo `json:"catInfo"`
	SuperLotID          int                     `json:"superlotID"`  // Found in React, not confirmed to be in json
	SuperLotQuantity    int                     `json:"superlotQty"` // Found in React, not confirmed to be in json
	DuplicateItems      []interface{}           `json:"duplicateItems"`
	NewLotIDLink        int                     `json:"newLotIDLink"`
	UsedLotIDLink       int                     `json:"usedLotIDLink"`
	ColorItems          []ColorItem             `json:"colorItems"`
	ItemRels            []interface{}           `json:"itemRels"`
	ItemRelLots         []interface{}           `json:"itemRelLots"`
	Recommended         []interface{}           `json:"recommended"`
	Wanted              Wanted                  `json:"wanted"`
	CartErrorMessage    string                  `json:"cartErrorMsg"`
	ReturnCode          int                     `json:"returnCode"`
	ReturnMessage       string                  `json:"returnMessage"`
	ErrorTicket         int                     `json:"errorTicket"`
	ProcessingTime      int                     `json:"procssingTime"`
}

type StoreItemCategoryInfo struct {
	CategoryID    int    `json:"catID"`
	CategoryName  string `json:"catName"`
	CategoryLevel int    `json:"catLevel"`
}

type LargeImage struct {
	ImageID        int         `json:"idImg"`
	TypeImageSmall interface{} `json:"typeImgS"`
	TypeImageLarge interface{} `json:"typeImgL"`
	URLSmall       string      `json:"strUrlS"`
	URLLarge       string      `json:"strUrlL"`
	ColorID        int         `json:"idColor"`
}

type ColorItem struct {
	InvID     int    `json:"invID"`
	ColorID   int    `json:"colorID"`
	ColorHex  string `json:"colorHex"`
	ColorName string `json:"colorName"`
}

type Wanted struct {
	Lists                []WantedListSummary `json:"lists"`
	WantedMoreName       string              `json:"wantedMoreName"`
	WantedRemarks        string              `json:"wantedRemarks"`
	WantedPrice          string              `json:"wantedPrice"`
	WantedQuantity       int                 `json:"wantedQty"`
	WantedQuantityFilled int                 `json:"wantedQtyFilled"`
	WantedQuantityHasAll bool                `json:"wantedQtyHasAll"`
	QuantityWarn         bool                `json:"qtyWarn"`
	PriceWarn            bool                `json:"priceWarn"`
	WantedListCount      int                 `json:"wantedListCnt"`
}

type WantedListSummary struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"desc"`
	Remarks        string `json:"remarks"`
	Condition      string `json:"condition"`
	Quantity       int    `json:"qty"`
	QuantityFilled int    `json:"qtyFilled"`
	QuantityWarn   bool   `json:"qtyWarn"`
	Price          string `json:"price"`
	PriceWarn      bool   `json:"priceWarn"`
}
