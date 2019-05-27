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
	TotalCount     int64         `json:"total_count"`
	ColorID        int64         `json:"idColor"`
	ResultsPerPage int64         `json:"rpp"`
	PageIndex      int64         `json:"pi"`
	List           []ProductList `json:"list"`
	ReturnCode     int64         `json:"returnCode"`
	ReturnMessage  string        `json:"returnMessage"`
	ErrorTicket    int64         `json:"errorTicket"`
	ProcessingTime int64         `json:"procssingTime"`
}

type ProductList struct {
	InvID                  int64        `json:"idInv"`
	Description            string       `json:"strDesc"`
	NewOrUsed              NewOrUsed    `json:"codeNew"`
	Completeness           Completeness `json:"codeComplete"`
	Quantity               int64        `json:"n4Qty"`
	ColorID                int64        `json:"idColor"`
	ColorIDDefault         int64        `json:"idColorDefault"`
	Color                  string       `json:"strColor"`
	ImageURL               string       `json:"strInvImgUrl"`           // URL of image provided by seller
	ImageID                int64        `json:"idInvImg"`               // ID of Image provided by seller
	ImageType              ImageType    `json:"typeInvImg"`             // Type of Image provided by seller
	ImageTypeDefault       ImageType    `json:"typeImgDefault"`         // Default catalog image type
	HasExtendedDescription int64        `json:"hasExtendedDescription"` // Has an extended description (0: false, 1: true)
	InstantCheckout        bool         `json:"instantCheckout"`
	SalePrice              string       `json:"mInvSalePrice"`     // Sale price in currency of store
	SalePriceDisplay       string       `json:"mDisplaySalePrice"` // Sale price in display currency of user
	SalePercent            int64        `json:"nSalePct"`          // Percent discounted
	Tier1Quantity          int64        `json:"nTier1Qty"`
	Tier2Quantity          int64        `json:"nTier2Qty"`
	Tier3Quantity          int64        `json:"nTier3Qty"`
	Tier1Price             string       `json:"nTier1InvPrice"`     // Price in currency of store
	Tier2Price             string       `json:"nTier2InvPrice"`     // Price in currency of store
	Tier3Price             string       `json:"nTier3InvPrice"`     // Price in currency of store
	Tier1PriceDisplay      string       `json:"nTier1DisplayPrice"` // Price in display currency of user
	Tier2PriceDisplay      string       `json:"nTier2DisplayPrice"` // Price in display currency of user
	Tier3PriceDisplay      string       `json:"nTier3DisplayPrice"` // Price in display currency of user
	Category               string       `json:"strCategory"`
	StoreName              string       `json:"strStorename"`
	StoreCurrencyID        int64        `json:"idCurrencyStore"`
	MinBuy                 string       `json:"mMinBuy"`
	SellerUsername         string       `json:"strSellerUsername"`
	SellerFeedbackScore    int64        `json:"n4SellerFeedbackScore"`
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
	url := fmt.Sprintf("%s/store/item.ajax?invID=%s&sid=%s&wantedMoreArrayID=%s", cloneBase, invID, storeID, wantedListArrayID)
	var storeItem StoreItem
	if err := c.doGet(url, &storeItem); err != nil {
		return nil, err
	}
	return &storeItem, checkResponse(storeItem.ReturnCode, storeItem.ReturnMessage)
}

type StoreItem struct {
	InvID               int64          `json:"invID"`
	Description         string         `json:"description"`
	ExtendedDescription string         `json:"extDescription"`
	InvQty              int64          `json:"invQty"`
	InvBulk             int64          `json:"invBulk"`
	InvSale             int64          `json:"invSale"`
	ItemType            string         `json:"itemType"`
	InvNew              string         `json:"invNew"`
	ColorID             int64          `json:"colorID"`
	ColorName           string         `json:"colorName"`
	ColorHex            string         `json:"colorHex"`
	InvDate             string         `json:"invDate"`
	ItemNo              string         `json:"itemNo"`
	ItemSeq             int64          `json:"itemSeq"`
	ItemID              int64          `json:"itemID"`
	ItemName            string         `json:"itemName"`
	ItemStatus          string         `json:"itemStatus"`
	ItemBrand           int64          `json:"itemBrand"`
	InvComplete         string         `json:"invComplete"`
	NativePrice         string         `json:"nativePrice"`
	SalePrice           string         `json:"salePrice"`
	InvPrice            string         `json:"invPrice"`
	InvTierQty1         int64          `json:"invTierQty1"`
	InvTierQty2         int64          `json:"invTierQty2"`
	InvTierQty3         int64          `json:"invTierQty3"`
	InvTierPrice1       string         `json:"invTierPrice1"`
	InvTierPrice2       string         `json:"invTierPrice2"`
	InvTierPrice3       string         `json:"invTierPrice3"`
	InvTierNativePrice1 string         `json:"invTierNativePrice1"`
	InvTierNativePrice2 string         `json:"invTierNativePrice2"`
	InvTierNativePrice3 string         `json:"invTierNativePrice3"`
	CartQty             int64          `json:"cartQty"`
	CartBindQty         int64          `json:"cartBindQty"`
	InvBindID           int64          `json:"invBindID"`
	InvBindQty          int64          `json:"invBindQty"`
	InvImageID          int64          `json:"invImgID"`
	InvImageType        string         `json:"invImgType"`
	InvURL              string         `json:"invURL"`
	ImageURL            string         `json:"imgURL"`
	LargeImages         []LargeImage   `json:"largeImgs"`
	CategoryInfo        []CategoryInfo `json:"catInfo"`
	DuplicateItems      []interface{}  `json:"duplicateItems"`
	NewLotIDLink        int64          `json:"newLotIDLink"`
	UsedLotIDLink       int64          `json:"usedLotIDLink"`
	ColorItems          []ColorItem    `json:"colorItems"`
	ItemRels            []interface{}  `json:"itemRels"`
	ItemRelLots         []interface{}  `json:"itemRelLots"`
	Recommended         []interface{}  `json:"recommended"`
	Wanted              Wanted         `json:"wanted"`
	ReturnCode          int64          `json:"returnCode"`
	ReturnMessage       string         `json:"returnMessage"`
	ErrorTicket         int64          `json:"errorTicket"`
	ProcessingTime      int64          `json:"procssingTime"`
}

type CategoryInfo struct {
	CategoryID    int64  `json:"catID"`
	CategoryName  string `json:"catName"`
	CategoryLevel int64  `json:"catLevel"`
}

type LargeImage struct {
	ImageID        int64       `json:"idImg"`
	TypeImageSmall interface{} `json:"typeImgS"`
	TypeImageLarge interface{} `json:"typeImgL"`
	URLSmall       string      `json:"strUrlS"`
	URLLarge       string      `json:"strUrlL"`
	ColorID        int64       `json:"idColor"`
}

type ColorItem struct {
	InvID     int64  `json:"invID"`
	ColorID   int64  `json:"colorID"`
	ColorHex  string `json:"colorHex"`
	ColorName string `json:"colorName"`
}

type Wanted struct {
	Lists                []WantedListSummary `json:"lists"`
	WantedMoreName       string              `json:"wantedMoreName"`
	WantedRemarks        string              `json:"wantedRemarks"`
	WantedPrice          string              `json:"wantedPrice"`
	WantedQuantity       int64               `json:"wantedQty"`
	WantedQuantityFilled int64               `json:"wantedQtyFilled"`
	WantedQuantityHasAll bool                `json:"wantedQtyHasAll"`
	QuantityWarn         bool                `json:"qtyWarn"`
	PriceWarn            bool                `json:"priceWarn"`
	WantedListCount      int64               `json:"wantedListCnt"`
}

type WantedListSummary struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"desc"`
	Remarks        string `json:"remarks"`
	Condition      string `json:"condition"`
	Quantity       int64  `json:"qty"`
	QuantityFilled int64  `json:"qtyFilled"`
	QuantityWarn   bool   `json:"qtyWarn"`
	Price          string `json:"price"`
	PriceWarn      bool   `json:"priceWarn"`
}
