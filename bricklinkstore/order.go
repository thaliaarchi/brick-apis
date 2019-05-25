package bricklinkstore

import (
	"fmt"
	"strings"
	"time"
)

// GetOrders retrieves a list of orders you received or placed. Direction can be "in" or "out".
func (c *Client) GetOrders(direction string) ([]Order, error) {
	url := fmt.Sprintf("/orders?direction=%s", direction)
	return c.getOrders(url)
}

// GetOrdersByStatus retrieves a list of orders you received or placed. Direction can be "in" or "out".
// Statuses can be provided to include or exclude orders in those statuses.
// Filed indicates whether the result retries filed or un-filed orders.
func (c *Client) GetOrdersByStatus(direction string, includeStatuses, excludeStatuses []string, filed bool) ([]Order, error) {
	statuses := includeStatuses
	if len(excludeStatuses) > 0 {
		statuses = append(statuses, "-"+strings.Join(excludeStatuses, ",-"))
	}
	status := strings.Join(statuses, ",")
	url := fmt.Sprintf("/orders?direction=%s&status=%s&filed=%t", direction, status, filed)
	return c.getOrders(url)
}

func (c *Client) getOrders(url string) ([]Order, error) {
	var orders ordersResponse
	if err := c.doGet(url, &orders); err != nil {
		return nil, err
	}
	return orders.Data, checkMeta(orders.Meta)
}

type ordersResponse struct {
	Meta meta    `json:"meta"`
	Data []Order `json:"data"`
}

// GetOrder retrieves the details of a specific order.
func (c *Client) GetOrder(id int64) (*Order, error) {
	url := fmt.Sprintf("/orders/%d", id)
	var order orderResponse
	if err := c.doGet(url, &order); err != nil {
		return nil, err
	}
	return &order.Data, checkMeta(order.Meta)
}

type orderResponse struct {
	Meta meta  `json:"meta"`
	Data Order `json:"data"`
}

// GetOrderItems retrieves a list of items for the specified order.
// Returns a list of batches, each containing a list of items.
func (c *Client) GetOrderItems(id int64) ([][]OrderItem, error) {
	url := fmt.Sprintf("/orders/%d/items", id)
	var items orderItemsResponse
	if err := c.doGet(url, &items); err != nil {
		return nil, err
	}
	return items.Data, checkMeta(items.Meta)
}

type orderItemsResponse struct {
	Meta meta          `json:"meta"`
	Data [][]OrderItem `json:"data"`
}

// Order contains details for an order.
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
	Payment           Payment     `json:"payment"`             // Information related to the payment of this order
	Shipping          Shipping    `json:"shipping"`            // Information related to the shipping. API name data normalization: http://apidev.bricklink.com/redmine/boards/1/topics/4
	Cost              Cost        `json:"cost"`                // Cost information for this order
	DisplayCost       Cost        `json:"disp_cost"`           // Cost information for this order in display currency
}

// Payment contains payment information for an order.
type Payment struct {
	Method       PaymentMethod `json:"method"`              // The payment method for this order
	CurrencyCode CurrencyCode  `json:"currency_code"`       // Currency code of the payment
	DatePaid     time.Time     `json:"date_paid,omitempty"` // The time the buyer paid
	Status       PaymentStatus `json:"status"`              // Payment status. Available statuses: https://www.bricklink.com/help.asp?helpID=121
}

// Shipping contains shipping and tracking information for an order.
type Shipping struct {
	Method          string    `json:"method"`        // Shipping method name
	MethodID        int64     `json:"method_id"`     // Shipping method ID
	TrackingNumbers string    `json:"tracking_no"`   // Tracking numbers for the shipping
	TrackingLink    string    `json:"tracking_link"` // URL for tracking the shipping. API-only field. It is not shown on the BrickLink pages
	DateShipped     time.Time `json:"date_shipped"`  // Shipping date. API-only field. It is not shown on the BrickLink pages
	Address         *Address  `json:"address"`       // The object representation of the shipping address
}

// Address contains a user's address. The split fields are given only if the user provided their address in normalized form.
type Address struct {
	Name        PersonName  `json:"name"`         // An object representation of a person's name
	Full        string      `json:"full"`         // The full address in not-well-formatted
	Address1    string      `json:"address1"`     // The first line of the address. It is provided only if a buyer updated his/her address and name as a normalized form
	Address2    string      `json:"address2"`     // The second line of the address. It is provided only if a buyer updated his/her address and name as a normalized form
	CountryCode CountryCode `json:"country_code"` // The country code
	City        string      `json:"city"`         // The city. It is provided only if a buyer updated his/her address and name as a normalized form
	State       string      `json:"state"`        // The state. It is provided only if a buyer updated his/her address and name as a normalized form
	PostalCode  string      `json:"postal_code"`  // The postal code. It is provided only if a buyer updated his/her address and name as a normalized form
}

// PersonName contains a user's name. The first and last names are given only if the user provided their name in normalized form.
type PersonName struct {
	Full  string `json:"full"`  // The full name of this person, including middle names, suffixes, etc.
	First string `json:"first"` // The family name (last name) of this person. It is provided only if a buyer updated his/her address and name as a normalized form
	Last  string `json:"last"`  // The given name (first name) of this person. It is provided only if a buyer updated his/her address and name as a normalized form
}

// Cost contains cost information for an order
type Cost struct {
	CurrencyCode CurrencyCode `json:"currency_code"`      // The currency code
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

// OrderItem is an item contained in an order
type OrderItem struct {
	InventoryID           int64        `json:"inventory_id"`                 // The ID of the inventory that includes the item
	Item                  CatalogItem  `json:"item"`                         // An object representation of the item
	ColorID               int64        `json:"color_id"`                     // The ID of the color of the item
	ColorName             string       `json:"color_name"`                   // Color name of the item
	Quantity              int64        `json:"quantity"`                     // The number of items purchased in this order
	NewOrUsed             NewOrUsed    `json:"new_or_used"`                  // Indicates whether the item is new or used (N: New, U: Used)
	Completeness          Completeness `json:"completeness,omitempty"`       // Indicates whether the set is complete or incomplete. This value is valid only for SET type. (C: Complete, B: Incomplete, S: Sealed)
	UnitPrice             float64      `json:"unit_price,string"`            // The original price of this item per sale unit
	UnitPriceFinal        float64      `json:"unit_price_final,string"`      // The unit price of this item after applying tiered pricing policy
	UnitPriceDisplay      float64      `json:"disp_unit_price,string"`       // The original price of this item per sale unit in display currency of the user
	UnitPriceFinalDisplay float64      `json:"disp_unit_price_final,string"` // The unit price of this item after applying tiered pricing policy in display currency of the user
	CurrencyCode          CurrencyCode `json:"currency_code"`                // The currency code of the price
	CurrencyCodeDisplay   CurrencyCode `json:"disp_currency_code"`           // The display currency code of the user
	Description           string       `json:"description"`                  // User remarks of the order item
	Remarks               string       `json:"remarks"`                      // User description of the order item
	Weight                float64      `json:"weight,string"`                // The weight of the item that overrides the catalog weight
	OrderCost             float64      `json:"order_cost,string"`
}

type OrderStatus string

const (
	OrderCompleted OrderStatus = "COMPLETED"
	OrderReceived  OrderStatus = "RECEIVED"
	OrderPending   OrderStatus = "PENDING"
	OrderPurged    OrderStatus = "PURGED"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusReceived  PaymentStatus = "Received"
	PaymentStatusNone      PaymentStatus = "None"
	PaymentStatusSent      PaymentStatus = "Sent"
)

type PaymentMethod string

const (
	PaymentMethodPayPal       PaymentMethod = "PayPal"
	PaymentMethodPayPalOnsite PaymentMethod = "PayPal (Onsite)"
)

// CountryCode is represented as ISO 3166-1 alpha-2 (exception: UK instead of GB).
// See: http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2.
type CountryCode string

// CurrencyCode is represented as ISO 4217. See: https://en.wikipedia.org/wiki/ISO_4217.
type CurrencyCode string

// BrickLink supports these currency codes. See: https://www.bricklink.com/help.asp?helpID=436.
const (
	CurrencyCodeARS CurrencyCode = "ARS" // Argentine Peso
	CurrencyCodeAUD CurrencyCode = "AUD" // Australian Dollar
	CurrencyCodeBRL CurrencyCode = "BRL" // Brazilian Real
	CurrencyCodeBGN CurrencyCode = "BGN" // Bulgarian Lev
	CurrencyCodeCAD CurrencyCode = "CAD" // Canadian Dollar
	CurrencyCodeCNY CurrencyCode = "CNY" // Chinese Yuan
	CurrencyCodeHRK CurrencyCode = "HRK" // Croatian Kuna
	CurrencyCodeCZK CurrencyCode = "CZK" // Czech Koruna
	CurrencyCodeDKK CurrencyCode = "DKK" // Danish Krone
	CurrencyCodeEUR CurrencyCode = "EUR" // Euro
	CurrencyCodeGTQ CurrencyCode = "GTQ" // Guatemalan Quetzal
	CurrencyCodeHKD CurrencyCode = "HKD" // Hong Kong Dollar
	CurrencyCodeHUF CurrencyCode = "HUF" // Hungarian Forint
	CurrencyCodeINR CurrencyCode = "INR" // Indian Rupee
	CurrencyCodeIDR CurrencyCode = "IDR" // Indonesian Rupiah
	CurrencyCodeILS CurrencyCode = "ILS" // Israeli New Shekel
	CurrencyCodeJPY CurrencyCode = "JPY" // Japanese Yen
	CurrencyCodeMOP CurrencyCode = "MOP" // Macau Pataca
	CurrencyCodeMYR CurrencyCode = "MYR" // Malaysian Ringgit
	CurrencyCodeMXN CurrencyCode = "MXN" // Mexican Peso
	CurrencyCodeNZD CurrencyCode = "NZD" // New Zealand Dollar
	CurrencyCodeNOK CurrencyCode = "NOK" // Norwegian Kroner
	CurrencyCodePHP CurrencyCode = "PHP" // Philippine Peso
	CurrencyCodePLN CurrencyCode = "PLN" // Polish Zloty
	CurrencyCodeGBP CurrencyCode = "GBP" // Pound Sterling
	CurrencyCodeRON CurrencyCode = "RON" // Romanian New Lei
	CurrencyCodeRUB CurrencyCode = "RUB" // Russian Rouble
	CurrencyCodeRSD CurrencyCode = "RSD" // Serbian Dinar
	CurrencyCodeSGD CurrencyCode = "SGD" // Singapore Dollar
	CurrencyCodeZAR CurrencyCode = "ZAR" // South African Rand
	CurrencyCodeKRW CurrencyCode = "KRW" // South Korean Won
	CurrencyCodeSEK CurrencyCode = "SEK" // Swedish Krona
	CurrencyCodeCHF CurrencyCode = "CHF" // Swiss Franc
	CurrencyCodeTWD CurrencyCode = "TWD" // Taiwan New Dollar
	CurrencyCodeTHB CurrencyCode = "THB" // Thai Baht
	CurrencyCodeTRY CurrencyCode = "TRY" // Turkish Lira
	CurrencyCodeUAH CurrencyCode = "UAH" // Ukraine Hryvnia
	CurrencyCodeUSD CurrencyCode = "USD" // US Dollar
)

type NewOrUsed string

const (
	NewOrUsedNew  NewOrUsed = "N"
	NewOrUsedUsed NewOrUsed = "U"
)

type Completeness string

const (
	CompletenessComplete   Completeness = "C"
	CompletenessIncomplete Completeness = "B"
	CompletenessSealed     Completeness = "S"
)
