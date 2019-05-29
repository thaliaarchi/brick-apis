package bricklinkstore

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GetItem returns information about the specified item in BrickLink catalog. CatalogItem contains all fields.
func (c *Client) GetItem(itemType ItemType, itemNo string) (*CatalogItem, error) {
	url := fmt.Sprintf("/items/%s/%s", itemType, itemNo)
	return c.getCatalogItem(url)
}

// GetItemImage returns image URL of the specified item by colors. CatalogItem includes No, Type, and ThumbnailURL.
func (c *Client) GetItemImage(itemType ItemType, itemNo string, colorID int) (*CatalogItem, error) {
	url := fmt.Sprintf("/items/%s/%s/images/%d", itemType, itemNo, colorID)
	return c.getCatalogItem(url)
}

func (c *Client) getCatalogItem(url string) (*CatalogItem, error) {
	var r catalogItemResponse
	if err := c.doGet(url, &r); err != nil {
		return nil, err
	}
	return &r.Data, checkMeta(r.Meta)
}

type catalogItemResponse struct {
	Meta meta        `json:"meta"`
	Data CatalogItem `json:"data"`
}

// GetSupersets returns a list of items that include any color of the specified item. CatalogItem includes No, Name, Type, and CategoryID.
func (c *Client) GetSupersets(itemType ItemType, itemNo string) ([]SupersetEntries, error) {
	url := fmt.Sprintf("/items/%s/%s/supersets", itemType, itemNo)
	return c.getSupersets(url)
}

// GetSupersetsByColor returns a list of items that include the specified item. CatalogItem includes No, Name, Type, and CategoryID.
func (c *Client) GetSupersetsByColor(itemType ItemType, itemNo string, colorID int) ([]SupersetEntries, error) {
	url := fmt.Sprintf("/items/%s/%s/supersets?color_id=%d", itemType, itemNo, colorID)
	return c.getSupersets(url)
}

func (c *Client) getSupersets(url string) ([]SupersetEntries, error) {
	var r supersetEntriesResponse
	if err := c.doGet(url, &r); err != nil {
		return nil, err
	}
	return r.Data, checkMeta(r.Meta)
}

type supersetEntriesResponse struct {
	Meta meta              `json:"meta"`
	Data []SupersetEntries `json:"data"`
}

// GetSubsets returns a list of items that are included in any color of the specified item. CatalogItem includes No, Name, Type, and CategoryID.
func (c *Client) GetSubsets(itemType ItemType, id string, includeBox, includeInstruction, breakMinifigs, breakSubsets bool) ([]SubsetEntries, error) {
	url := fmt.Sprintf("/items/%s/%s/subsets?%s", itemType, id, subsetParams(includeBox, includeInstruction, breakMinifigs, breakSubsets))
	return c.getSubsets(url)
}

// GetSubsetsByColor returns a list of items that are included in the specified item. CatalogItem includes No, Name, Type, and CategoryID.
func (c *Client) GetSubsetsByColor(itemType ItemType, id string, colorID int, includeBox, includeInstruction, breakMinifigs, breakSubsets bool) ([]SubsetEntries, error) {
	url := fmt.Sprintf("/items/%s/%s/subsets?color_id=%d&%s", itemType, id, colorID, subsetParams(includeBox, includeInstruction, breakMinifigs, breakSubsets))
	return c.getSubsets(url)
}

func (c *Client) getSubsets(url string) ([]SubsetEntries, error) {
	var r subsetEntriesResponse
	if err := c.doGet(url, &r); err != nil {
		return nil, err
	}
	return r.Data, checkMeta(r.Meta)
}

func subsetParams(includeBox, includeInstruction, breakMinifigs, breakSubsets bool) string {
	return fmt.Sprintf("box=%t&instruction=%t&break_minifigs=%t&break_subsets=%t", includeBox, includeInstruction, breakMinifigs, breakSubsets)
}

type subsetEntriesResponse struct {
	Meta meta            `json:"meta"`
	Data []SubsetEntries `json:"data"`
}

// GetPriceGuide returns the price statistics of the specified item in BrickLink catalog. CatalogItem includes No and Type.
func (c *Client) GetPriceGuide(itemType ItemType, itemNo string, options *PriceGuideOptions) (*PriceGuide, error) {
	url := fmt.Sprintf("/items/%s/%s/price%s", itemType, itemNo, toParams(options))
	var r priceGuideResponse
	if err := c.doGet(url, &r); err != nil {
		return nil, err
	}
	return &r.Data, checkMeta(r.Meta)
}

type priceGuideResponse struct {
	Meta meta       `json:"meta"`
	Data PriceGuide `json:"data"`
}

// PriceGuideOptions contains optional parameters for GetPriceGuide
type PriceGuideOptions struct {
	ColorID      int // The color of the item or 0 for all colors
	GuideType    GuideType
	NewOrUsed    NewOrUsed
	CountryCode  CountryCode
	Region       Region
	CurrencyCode CurrencyCode
	VAT          IncludeVAT
}

func toParams(o *PriceGuideOptions) string {
	if o == nil {
		return ""
	}
	var params url.Values
	if o.ColorID != 0 {
		params.Set("color_id", strconv.Itoa(o.ColorID))
	}
	if o.GuideType != "" {
		params.Set("guide_type", string(o.GuideType))
	}
	if o.NewOrUsed != "" {
		params.Set("new_or_used", string(o.NewOrUsed))
	}
	if o.CountryCode != "" {
		params.Set("country_code", string(o.CountryCode))
	}
	if o.Region != "" {
		params.Set("region", string(o.Region))
	}
	if o.CurrencyCode != "" {
		params.Set("currency_code", string(o.CurrencyCode))
	}
	if o.VAT != "" {
		params.Set("vat", string(o.VAT))
	}
	if str := params.Encode(); str != "" {
		return "?" + str
	}
	return ""
}

// GetKnownColors returns currently known colors of the item.
func (c *Client) GetKnownColors(itemType ItemType, id string) ([]KnownColor, error) {
	url := fmt.Sprintf("/items/%s/%s/colors", itemType, id)
	var r knownColorsResponse
	if err := c.doGet(url, &r); err != nil {
		return nil, err
	}
	return r.Data, checkMeta(r.Meta)
}

type knownColorsResponse struct {
	Meta meta         `json:"meta"`
	Data []KnownColor `json:"data"`
}

type CatalogItem struct {
	No           string   `json:"no"`                      // Item's identification number in BrickLink catalog
	Name         string   `json:"name"`                    // The name of the item
	Type         ItemType `json:"type"`                    // The type of the item (MINIFIG, PART, SET, BOOK, GEAR, CATALOG, INSTRUCTION, UNSORTED_LOT, ORIGINAL_BOX)
	CategoryID   int      `json:"category_id"`             // The main category of the item
	AlternateNo  string   `json:"alternate_no,omitempty"`  // Alternate item number. Alternate item number: https://www.bricklink.com/help.asp?helpID=599
	ImageURL     string   `json:"image_url"`               // Image link for this item
	ThumbnailURL string   `json:"thumbnail_url"`           // Image thumbnail link for this item
	Weight       float64  `json:"weight,string"`           // The weight of the item in grams (with 2 decimal places)
	DimX         float64  `json:"dim_x,string"`            // Length of the item. Item dimensions with 2 decimal places: https://www.bricklink.com/help.asp?helpID=261
	DimY         float64  `json:"dim_y,string"`            // Width of the item. Item dimensions with 2 decimal places: https://www.bricklink.com/help.asp?helpID=261
	DimZ         float64  `json:"dim_z,string"`            // Height of the item. Item dimensions with 2 decimal places: https://www.bricklink.com/help.asp?helpID=261
	YearReleased int      `json:"year_released"`           // Item year of release. https://www.bricklink.com/help.asp?helpID=2004
	Description  string   `json:"description,omitempty"`   // Short description for this item
	IsObsolete   bool     `json:"is_obsolete"`             // Indicates whether the item is obsolete
	LanguageCode string   `json:"language_code,omitempty"` // Item language code
}

type SupersetEntries struct {
	ColorID int             `json:"color_id"` // The ID of the color of the item
	Entries []SupersetEntry `json:"entries"`  // A list of the items that include the specified item
}

type SupersetEntry struct {
	Item      CatalogItem `json:"item"`       // An object representation of the super item that includes the specified item
	Quantity  int         `json:"quantity"`   // Indicates that how many specified items are included in this super item
	AppearsAs AppearsAs   `json:"appears_as"` // Indicates how an entry in an inventory appears as (A: Alternate, C: Counterpart, E: Extra, R: Regular)
}

type SubsetEntries struct {
	MatchNo int           `json:"match_no"` // An identification number given to a matching group that consists of regular items and alternate items. 0 if there is no matching of alternative item
	Entries []SubsetEntry `json:"entries"`  // A list of the items included in the specified item
}

type SubsetEntry struct {
	Item          CatalogItem `json:"item"`           // An object representation of the item that is included in the specified item
	ColorID       int         `json:"color_id"`       // The ID of the color of the item
	Quantity      int         `json:"quantity"`       // The number of items that are included in
	ExtraQuantity int         `json:"extra_quantity"` // The number of items that are appear as "extra" item
	IsAlternate   bool        `json:"is_alternate"`   // Indicates that the item is appear as "alternate" item in this specified item
	IsCounterpart bool        `json:"is_counterpart"`
}

type PriceGuide struct {
	Item          CatalogItem   `json:"item"`                 // An object representation of the item
	NewOrUsed     string        `json:"new_or_used"`          // Indicates whether the price guide is for new or used (N: New, U: Used)
	CurrencyCode  string        `json:"currency_code"`        // The currency code of the price
	MinPrice      float64       `json:"min_price,string"`     // The lowest price of the item (in stock / that was sold for last 6 months)
	MaxPrice      float64       `json:"max_price,string"`     // The highest price of the item (in stock / that was sold for last 6 months)
	AvgPrice      float64       `json:"avg_price,string"`     // The average price of the item (in stock / that was sold for last 6 months)
	QtyAvgPrice   float64       `json:"qty_avg_price,string"` // The average price of the item (in stock / that was sold for last 6 months) by quantity
	UnitQuantity  int           `json:"unit_quantity"`        // The number of inventories that include the item / The number of times the item has been sold for last 6 months
	TotalQuantity int           `json:"total_quantity"`       // The total number of the items in stock / The number of items has been sold for last 6 months
	PriceDetail   []PriceDetail `json:"price_detail"`         // A list of objects that represent the detailed information of the price
}

type PriceDetail struct {
	Quantity           int         `json:"quantity"`            // The number of the items in the inventory
	QuantityDeprecated int         `json:"qunatity"`            // Deprecated, typo
	UnitPrice          float64     `json:"unit_price,string"`   // The original price of this item per sale unit
	ShippingAvailable  bool        `json:"shipping_available"`  // Indicates whether or not the seller ships to your country (based on the user profile). Only included for in stock
	SellerCountryCode  CountryCode `json:"seller_country_code"` // The country code of the seller's location. Only included for last 6 months.
	BuyerCountryCode   CountryCode `json:"buyer_country_code"`  // The country code of the buyer's location. Only included for last 6 months.
	DateOrdered        time.Time   `json:"date_ordered"`        // The time the order was created. Only included for last 6 months.
}

type KnownColor struct {
	ColorID  int `json:"color_id"` // Color ID
	Quantity int `json:"quantity"` // The quantity of items in that color
}

type ItemType string

const (
	ItemTypeMinifig     ItemType = "MINIFIG"
	ItemTypePart        ItemType = "PART"
	ItemTypeSet         ItemType = "SET"
	ItemTypeBook        ItemType = "BOOK"
	ItemTypeGear        ItemType = "GEAR"
	ItemTypeCatalog     ItemType = "CATALOG"
	ItemTypeInstruction ItemType = "INSTRUCTION"
	ItemTypeUnsortedLot ItemType = "UNSORTED_LOT"
	ItemTypeOriginalBox ItemType = "ORIGINAL_BOX"
)

type AppearsAs string

const (
	AppearsAsAlternate   AppearsAs = "A"
	AppearsAsCounterpart AppearsAs = "C"
	AppearsAsExtra       AppearsAs = "E"
	AppearsAsRegular     AppearsAs = "R"
)

// GuideType indicates which price guide statistics to be provided.
type GuideType string

// Available values for GuideType. See: http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki/CatalogMethod#-Parameters-5
const (
	GuideTypeStock GuideType = "stock" // Current Items for Sale (default)
	GuideTypeSold  GuideType = "sold"  // Last 6 Months Sales
)

// Region is a geographical area for store grouping.
type Region string

// Available values for Region. See: http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki/CatalogMethod#-Parameters-5
const (
	RegionAfrica       Region = "africa"        // Africa
	RegionAsia         Region = "asia"          // Asia
	RegionEU           Region = "eu"            // European Union
	RegionEurope       Region = "europe"        // Europe
	RegionMiddleEast   Region = "middle_east"   // Middle East
	RegionNorthAmerica Region = "north_america" // North America
	RegionOceania      Region = "oceania"       // Australia & Oceania
	RegionSouthAmerica Region = "south_america" // South America
)

// IncludeVAT indicates that price will include VAT for the items of VAT enabled stores.
type IncludeVAT string

// Available values for IncludeVAT. See: http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki/CatalogMethod#-Parameters-5.
const (
	IncludeVATNo     IncludeVAT = "N" // Exclude VAT (default)
	IncludeVATYes    IncludeVAT = "Y" // Include VAT
	IncludeVATNorway IncludeVAT = "O" // Include VAT as Norway settings
)
