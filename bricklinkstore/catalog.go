package bricklinkstore

import "fmt"

// GetItem returns information about the specified item in BrickLink catalog. CatalogItem contains all fields.
func (c *Client) GetItem(itemType ItemType, id string) (*CatalogItem, error) {
	url := fmt.Sprintf("/items/%s/%s", itemType, id)
	return c.getCatalogItem(url)
}

// GetItemImage returns image URL of the specified item by colors. CatalogItem includes No, Type, and ThumbnailURL.
func (c *Client) GetItemImage(itemType ItemType, id string, colorID int) (*CatalogItem, error) {
	url := fmt.Sprintf("/items/%s/%s/images/%d", itemType, id, colorID)
	return c.getCatalogItem(url)
}

func (c *Client) getCatalogItem(url string) (*CatalogItem, error) {
	var item catalogItemResponse
	if err := c.doGet(url, &item); err != nil {
		return nil, err
	}
	return &item.Data, checkMeta(item.Meta)
}

type catalogItemResponse struct {
	Meta meta        `json:"meta"`
	Data CatalogItem `json:"data"`
}

// GetSupersets returns a list of items that include any color of the specified item. CatalogItem includes No, Name, Type, and CategoryID.
func (c *Client) GetSupersets(itemType ItemType, id string) ([]SupersetEntries, error) {
	url := fmt.Sprintf("/items/%s/%s/supersets", itemType, id)
	return c.getSupersets(url)
}

// GetSupersetsByColor returns a list of items that include the specified item. CatalogItem includes No, Name, Type, and CategoryID.
func (c *Client) GetSupersetsByColor(itemType ItemType, id string, colorID int) ([]SupersetEntries, error) {
	url := fmt.Sprintf("/items/%s/%s/supersets?color_id=%d", itemType, id, colorID)
	return c.getSupersets(url)
}

func (c *Client) getSupersets(url string) ([]SupersetEntries, error) {
	var supersets supersetEntriesResponse
	if err := c.doGet(url, &supersets); err != nil {
		return nil, err
	}
	return supersets.Data, checkMeta(supersets.Meta)
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
	var subsets subsetEntriesResponse
	if err := c.doGetAndSave(url, &subsets, "subsets.json"); err != nil {
		return nil, err
	}
	return subsets.Data, checkMeta(subsets.Meta)
}

func subsetParams(includeBox, includeInstruction, breakMinifigs, breakSubsets bool) string {
	return fmt.Sprintf("box=%t&instruction=%t&break_minifigs=%t&break_subsets=%t", includeBox, includeInstruction, breakMinifigs, breakSubsets)
}

type subsetEntriesResponse struct {
	Meta meta            `json:"meta"`
	Data []SubsetEntries `json:"data"`
}

type CatalogItem struct {
	No           string   `json:"no"`                      // Item's identification number in BrickLink catalog
	Name         string   `json:"name"`                    // The name of the item
	Type         ItemType `json:"type"`                    // The type of the item (MINIFIG, PART, SET, BOOK, GEAR, CATALOG, INSTRUCTION, UNSORTED_LOT, ORIGINAL_BOX)
	CategoryID   int64    `json:"category_id"`             // The main category of the item
	AlternateNo  string   `json:"alternate_no,omitempty"`  // Alternate item number. Alternate item number: https://www.bricklink.com/help.asp?helpID=599
	ImageURL     string   `json:"image_url"`               // Image link for this item
	ThumbnailURL string   `json:"thumbnail_url"`           // Image thumbnail link for this item
	Weight       float64  `json:"weight,string"`           // The weight of the item in grams (with 2 decimal places)
	DimX         float64  `json:"dim_x,string"`            // Length of the item. Item dimensions with 2 decimal places: https://www.bricklink.com/help.asp?helpID=261
	DimY         float64  `json:"dim_y,string"`            // Width of the item. Item dimensions with 2 decimal places: https://www.bricklink.com/help.asp?helpID=261
	DimZ         float64  `json:"dim_z,string"`            // Height of the item. Item dimensions with 2 decimal places: https://www.bricklink.com/help.asp?helpID=261
	YearReleased int64    `json:"year_released"`           // Item year of release. https://www.bricklink.com/help.asp?helpID=2004
	Description  string   `json:"description,omitempty"`   // Short description for this item
	IsObsolete   bool     `json:"is_obsolete"`             // Indicates whether the item is obsolete
	LanguageCode string   `json:"language_code,omitempty"` // Item language code
}

type SupersetEntries struct {
	ColorID int64           `json:"color_id"` // The ID of the color of the item
	Entries []SupersetEntry `json:"entries"`  // A list of the items that include the specified item
}

type SupersetEntry struct {
	Item      CatalogItem `json:"item"`       // An object representation of the super item that includes the specified item
	Quantity  int64       `json:"quantity"`   // Indicates that how many specified items are included in this super item
	AppearsAs AppearsAs   `json:"appears_as"` // Indicates how an entry in an inventory appears as (A: Alternate, C: Counterpart, E: Extra, R: Regular)
}

type SubsetEntries struct {
	MatchNo int64         `json:"match_no"` // An identification number given to a matching group that consists of regular items and alternate items. 0 if there is no matching of alternative item
	Entries []SubsetEntry `json:"entries"`  // A list of the items included in the specified item
}

type SubsetEntry struct {
	Item          CatalogItem `json:"item"`           // An object representation of the item that is included in the specified item
	ColorID       int64       `json:"color_id"`       // The ID of the color of the item
	Quantity      int64       `json:"quantity"`       // The number of items that are included in
	ExtraQuantity int64       `json:"extra_quantity"` // The number of items that are appear as "extra" item
	IsAlternate   bool        `json:"is_alternate"`   // Indicates that the item is appear as "alternate" item in this specified item
	IsCounterpart bool        `json:"is_counterpart"`
}

type ItemType string
type AppearsAs string

const (
	ItemTypeMinifig      ItemType  = "MINIFIG"
	ItemTypePart         ItemType  = "PART"
	ItemTypeSet          ItemType  = "SET"
	ItemTypeBook         ItemType  = "BOOK"
	ItemTypeGear         ItemType  = "GEAR"
	ItemTypeCatalog      ItemType  = "CATALOG"
	ItemTypeInstruction  ItemType  = "INSTRUCTION"
	ItemTypeUnsortedLot  ItemType  = "UNSORTED_LOT"
	ItemTypeOriginalBox  ItemType  = "ORIGINAL_BOX"
	AppearsAsAlternate   AppearsAs = "A"
	AppearsAsCounterpart AppearsAs = "C"
	AppearsAsExtra       AppearsAs = "E"
	AppearsAsRegular     AppearsAs = "R"
)
