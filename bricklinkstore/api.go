package bricklinkstore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andrewarchi/bricklink-buy/credentials"
	"github.com/mrjones/oauth"
)

const (
	apiBase = "https://api.bricklink.com/api/store/v1"
)

type Client struct {
	client *http.Client
}

func NewClient(cred *credentials.BrickLinkStore) (*Client, error) {
	consumer := oauth.NewConsumer(cred.ConsumerKey, cred.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cred.Token, Secret: cred.TokenSecret}
	client, err := consumer.MakeHttpClient(accessToken)
	return &Client{client}, err
}

func (c *Client) GetOrderList() ([]Order, error) {
	url := apiBase + "/orders?direction=out"
	var orders OrderListResponse
	if err := c.doGet(url, &orders); err != nil {
		return nil, err
	}
	return orders.Data, nil
}

func (c *Client) GetOrder(id int64) (*Order, error) {
	url := fmt.Sprintf(apiBase+"/orders/%d", id)
	var order OrderResponse
	if err := c.doGet(url, &order); err != nil {
		return nil, err
	}
	return &order.Data, nil
}

func (c *Client) GetOrderItems(id int64) ([][]OrderItem, error) {
	url := fmt.Sprintf(apiBase+"/orders/%d/items", id)
	var items OrderItemsResponse
	if err := c.doGet(url, &items); err != nil {
		return nil, err
	}
	return items.Data, nil
}

func (c *Client) GetColorList() ([]Color, error) {
	url := apiBase + "/colors"
	var colors ColorListResponse
	if err := c.doGet(url, &colors); err != nil {
		return nil, err
	}
	return colors.Data, nil
}

func (c *Client) GetColor(id int64) (*Color, error) {
	url := fmt.Sprintf(apiBase+"/colors/%d", id)
	var color ColorResponse
	if err := c.doGet(url, &color); err != nil {
		return nil, err
	}
	return &color.Data, nil
}

func (c *Client) doGet(url string, v interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (o *Order) printUnknownValues() {
	switch o.Status {
	default:
		fmt.Printf("OrderStatus: %s", o.Status)
	case OrderCompleted:
	case OrderReceived:
	case OrderPurged:
	}
	if o.Payment != nil {
		switch o.Payment.Method {
		default:
			fmt.Printf("PaymentMethod: %s", o.Payment.Method)
		case PaymentPayPal, PaymentPayPalOnsite:
		}
		switch o.Payment.CurrencyCode {
		default:
			fmt.Printf("CurrencyCode: %s", o.Payment.CurrencyCode)
		case CurrencyCad, CurrencyEur, CurrencyHuf, CurrencyUsd:
		}
		switch o.Payment.Status {
		default:
			fmt.Printf("PaymentStatus: %s", o.Payment.Status)
		case PaymentCompleted, PaymentReceived, PaymentNone, PaymentSent:
		}
	}
}

type OrderListResponse struct {
	Meta Meta    `json:"meta"`
	Data []Order `json:"data"`
}

type OrderResponse struct {
	Meta Meta  `json:"meta"`
	Data Order `json:"data"`
}

type OrderItemsResponse struct {
	Meta Meta          `json:"meta"`
	Data [][]OrderItem `json:"data"`
}

type ColorListResponse struct {
	Meta Meta    `json:"meta"`
	Data []Color `json:"data"`
}

type ColorResponse struct {
	Meta Meta  `json:"meta"`
	Data Color `json:"data"`
}

type Meta struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Code        int64  `json:"code"`
}

type Order struct {
	OrderID           int64       `json:"order_id"`            // Unique identifier for this order for internal use
	DateOrdered       time.Time   `json:"date_ordered"`        // The time the order was created
	DateStatusChanged time.Time   `json:"date_status_changed"` // The time the order status was last modified
	SellerName        string      `json:"seller_name"`         // The username of the seller in BL
	StoreName         string      `json:"store_name"`          // The store name displayed on BL store pages
	BuyerName         string      `json:"buyer_name"`          // The username of the buyer in BL
	BuyerEmail        string      `json:"buyer_email"`         // E-mail address of the buyer
	BuyerOrderCount   int64       `json:"buyer_order_count"`   // Total count of all orders placed by the buyer in the seller's store. Includes the order just placed and also purged orders
	RequireInsurance  bool        `json:"require_insurance"`   // Indicates whether the buyer requests insurance for this order
	Status            OrderStatus `json:"status"`              // The status of an order. Available statuses: http://www.bricklink.com/help.asp?helpID=41
	IsInvoiced        bool        `json:"is_invoiced"`         // Indicates whether the order invoiced
	IsFiled           bool        `json:"is_filed"`            // Indicates whether the order is filed
	DriveThruSent     bool        `json:"drive_thru_sent"`     // Indicates whether "Thank You, Drive Thru!" email has been sent
	Remarks           string      `json:"remarks,omitempty"`   // User remarks for this order
	TotalCount        int64       `json:"total_count"`         // The total number of items included in this order
	UniqueCount       int64       `json:"unique_count"`        // The unique number of items included in this order
	TotalWeight       float64     `json:"total_weight,string"` // The total weight of the items ordered. It applies the seller's custom weight when present to override the catalog weight. 0 if the order includes at least one item without any weight information or incomplete set
	Payment           *Payment    `json:"payment"`             // Information related to the payment of this order
	Shipping          *Shipping   `json:"shipping"`            // Information related to the shipping. API name data normalization: http://apidev.bricklink.com/redmine/boards/1/topics/4
	Cost              *Cost       `json:"cost"`                // Cost information for this order
	DisplayCost       *Cost       `json:"disp_cost"`           // Cost information for this order in display currency
}

type Payment struct {
	Method       PaymentMethod `json:"method"`              // The payment method for this order
	CurrencyCode CurrencyCode  `json:"currency_code"`       // Currency code of the payment. ISO 4217: http://en.wikipedia.org/wiki/ISO_4217
	DatePaid     time.Time     `json:"date_paid,omitempty"` // The time the buyer paid
	Status       PaymentStatus `json:"status"`              // Payment status. Available statuses: https://www.bricklink.com/help.asp?helpID=121
}

type Shipping struct {
	Method          string    `json:"method"`        // Shipping method name
	MethodID        int64     `json:"method_id"`     // Shipping method ID
	TrackingNumbers string    `json:"tracking_no"`   // Tracking numbers for the shipping
	TrackingLink    string    `json:"tracking_link"` // URL for tracking the shipping. API-only field. It is not shown on the BrickLink pages
	DateShipped     time.Time `json:"date_shipped"`  // Shipping date. API-only field. It is not shown on the BrickLink pages
	Address         *Address  `json:"address"`       // The object representation of the shipping address
}

type Address struct {
	Name        *PersonName `json:"name"`         // An object representation of a person's name
	Full        string      `json:"full"`         // The full address in not-well-formatted
	Address1    string      `json:"address1"`     // The first line of the address. It is provided only if a buyer updated his/her address and name as a normalized form
	Address2    string      `json:"address2"`     // The second line of the address. It is provided only if a buyer updated his/her address and name as a normalized form
	CountryCode string      `json:"country_code"` // The country code. ISO 3166-1 alpha-2 (exception: UK instead of GB) http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	City        string      `json:"city"`         // The city. It is provided only if a buyer updated his/her address and name as a normalized form
	State       string      `json:"state"`        // The state. It is provided only if a buyer updated his/her address and name as a normalized form
	PostalCode  string      `json:"postal_code"`  // The postal code. It is provided only if a buyer updated his/her address and name as a normalized form
}

type PersonName struct {
	Full  string `json:"full"`  // The full name of this person, including middle names, suffixes, etc.
	First string `json:"first"` // The family name (last name) of this person. It is provided only if a buyer updated his/her address and name as a normalized form
	Last  string `json:"last"`  // The given name (first name) of this person. It is provided only if a buyer updated his/her address and name as a normalized form
}

type Cost struct {
	CurrencyCode CurrencyCode `json:"currency_code"`      // The currency code. ISO 4217: http://en.wikipedia.org/wiki/ISO_4217
	Subtotal     float64      `json:"subtotal,string"`    // The total price for the order exclusive of shipping and other costs. This must equal the sum of all the items
	GrandTotal   float64      `json:"grand_total,string"` // The total price for the order inclusive of tax, shipping and other costs
	Etc1         float64      `json:"etc1,string"`        // Extra charge for this order (tax, packing, etc.)
	Etc2         float64      `json:"etc2,string"`        // Extra charge for this order (tax, packing, etc.)
	Insurance    float64      `json:"insurance,string"`   // Insurance cost
	Shipping     float64      `json:"shipping,string"`    // Shipping cost
	Credit       float64      `json:"credit,string"`      // Credit applied to this order
	Coupon       float64      `json:"coupon,string"`      // Amount of coupon discount
	SalesTax     float64      `json:"salesTax,string"`
	VATRate      float64      `json:"vat_rate,string"`   // VAT percentage applied to this order
	VATAmount    float64      `json:"vat_amount,string"` // Total amount of VAT included in the grand_total price
}

// http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki/ColorResource
type Color struct {
	ColorID   int64     `json:"color_id"`   // ID of the color
	ColorName string    `json:"color_name"` // The name of the color
	ColorCode string    `json:"color_code"` // HTML color code of this color
	ColorType ColorType `json:"color_type"` // The name of the color group to which this color belongs
}

type OrderItem struct {
	InventoryID        int64        `json:"inventory_id"`
	Item               Item         `json:"item"`
	ColorID            int64        `json:"color_id"`
	ColorName          string       `json:"color_name"`
	Quantity           int64        `json:"quantity"`
	NewOrUsed          NewOrUsed    `json:"new_or_used"`
	Completeness       Completeness `json:"completeness,omitempty"`
	UnitPrice          string       `json:"unit_price"`
	Description        string       `json:"description"`
	Remarks            string       `json:"remarks"`
	DispUnitPrice      string       `json:"disp_unit_price"`
	DispUnitPriceFinal string       `json:"disp_unit_price_final"`
	UnitPriceFinal     string       `json:"unit_price_final"`
	OrderCost          string       `json:"order_cost"`
	CurrencyCode       CurrencyCode `json:"currency_code"`
	DispCurrencyCode   CurrencyCode `json:"disp_currency_code"`
	Weight             string       `json:"weight"`
}

type Item struct {
	ItemNumber string   `json:"no"`
	Name       string   `json:"name"`
	Type       ItemType `json:"type"`
	CategoryID int64    `json:"category_id"`
}

type OrderStatus string
type PaymentStatus string
type PaymentMethod string
type CurrencyCode string
type ColorType string
type ItemType string
type NewOrUsed string
type Completeness string

const (
	OrderCompleted       OrderStatus   = "COMPLETED"
	OrderReceived        OrderStatus   = "RECEIVED"
	OrderPurged          OrderStatus   = "PURGED"
	PaymentCompleted     PaymentStatus = "Completed"
	PaymentReceived      PaymentStatus = "Received"
	PaymentNone          PaymentStatus = "None"
	PaymentSent          PaymentStatus = "Sent"
	PaymentPayPal        PaymentMethod = "PayPal"
	PaymentPayPalOnsite  PaymentMethod = "PayPal (Onsite)"
	CurrencyCad          CurrencyCode  = "CAD"
	CurrencyEur          CurrencyCode  = "EUR"
	CurrencyHuf          CurrencyCode  = "HUF"
	CurrencyUsd          CurrencyCode  = "USD"
	ColorTypeBrickArms   ColorType     = "BrickArms"
	ColorTypeChrome      ColorType     = "Chrome"
	ColorTypeGlitter     ColorType     = "Glitter"
	ColorTypeMetallic    ColorType     = "Metallic"
	ColorTypeMilky       ColorType     = "Milky"
	ColorTypeModulex     ColorType     = "Modulex"
	ColorTypePearl       ColorType     = "Pearl"
	ColorTypeSolid       ColorType     = "Solid"
	ColorTypeSpeckle     ColorType     = "Speckle"
	ColorTypeTransparent ColorType     = "Transparent"
	ItemTypeMinifig      ItemType      = "MINIFIG"
	ItemTypePart         ItemType      = "PART"
	ItemTypeSet          ItemType      = "SET"
	ItemTypeBook         ItemType      = "BOOK"
	ItemTypeGear         ItemType      = "GEAR"
	ItemTypeCatalog      ItemType      = "CATALOG"
	ItemTypeInstruction  ItemType      = "INSTRUCTION"
	ItemTypeUnsortedLot  ItemType      = "UNSORTED_LOT"
	ItemTypeOriginalBox  ItemType      = "ORIGINAL_BOX"
	NewOrUsedNew         NewOrUsed     = "N"
	CompletenessNA       Completeness  = ""
	CompletenessComplete Completeness  = "C"
)
