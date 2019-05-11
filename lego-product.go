package main

type ProductInformation struct {
	Product                Product     `json:"Product"`
	Bricks                 []Brick     `json:"Bricks"`
	ImageBaseURL           string      `json:"ImageBaseUrl"`
	UnAvailableInformation interface{} `json:"UnAvailableInformation"`
}

type Brick struct {
	ItemNo            int64              `json:"ItemNo"`
	ItemDescr         string             `json:"ItemDescr"`
	ColourLikeDescr   string             `json:"ColourLikeDescr"`
	ColourDescr       string             `json:"ColourDescr"`
	MaingroupDescr    string             `json:"MaingroupDescr"`
	Asset             string             `json:"Asset"`
	MaxQty            int64              `json:"MaxQty"`
	IP                bool               `json:"Ip"`
	Price             float64            `json:"Price"`
	CID               string             `json:"CId"`
	SQty              int64              `json:"SQty"`
	DesignID          int64              `json:"DesignId"`
	PriceStr          string             `json:"PriceStr"`
	PriceWithTaxStr   string             `json:"PriceWithTaxStr"`
	ItemUnavailable   bool               `json:"ItemUnavailable"`
	UnavailableLink   *UnavailableLink   `json:"UnavailableLink"`
	UnavailableReason *UnavailableReason `json:"UnavailableReason"`
}

type UnavailableLink struct {
	URL   string `json:"Url"`
	Title string `json:"Title"`
}

type UnavailableReason struct {
	ReasonText        string        `json:"ReasonText"`
	LinkText          string        `json:"LinkText"`
	ID                string        `json:"ID"`
	Path              string        `json:"Path"`
	Key               string        `json:"Key"`
	HasVersion        bool          `json:"HasVersion"`
	RestrictedMarkets []interface{} `json:"RestrictedMarkets"`
}

type Product struct {
	ProductNo   string `json:"ProductNo"`
	ProductName string `json:"ProductName"`
	ItemNo      string `json:"ItemNo"`
	Asset       string `json:"Asset"`
}

/*
CId:
""
"USD"

ColourLikeDescr:
""
"Black"
"Blue"
"Grey"
"Purple"
"Red"
"White"
"Yellow"
*/

/*
MaingroupDescr:
"Bricks"
"Bricks, Special"
"Bricks, Special Circles And Angles"
*/

/*
ColorDescr:
"BLACK"
"BR.BLUE"
"BR.BLUEGREEN"
"BR.GREEN"
"BR.ORANGE"
"BR.RED"
"BR.YEL"
"BRICK-YEL"
"DK. BROWN"
"DK. ST. GREY"
"DK.GREEN"
"EARTH BLUE"
"FL. YELL-ORA"
"GOLD INK"
"LGH. PURPLE"
"M. LILAC"
"M. NOUGAT"
"MD.BLUE"
"MED. ST-GREY"
"MEDIUM AZUR"
"NEW DARK RED"
"OLIVE GREEN"
"RED. BROWN"
"SAND YELLOW"
"SILVER INK"
"SILVER MET."
"TR. BR. ORANGE"
"TR."
"TR.GREEN"
"TR.L.BLUE"
"TR.RED"
"W.GOLD"
"WHITE"
*/
