package main

import (
	"fmt"
	"net/http"

	"github.com/mrjones/oauth"
)

const (
	apiBase = "https://api.bricklink.com/api/store/v1"
)

func createBLStoreClient(cred *credentials) (*http.Client, error) {
	consumer := oauth.NewConsumer(cred.ConsumerKey, cred.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cred.Token, Secret: cred.TokenSecret}
	return consumer.MakeHttpClient(accessToken)
}

func getOrderList(client *http.Client) (*OrderListResponse, error) {
	url := apiBase + "/orders?direction=out"
	var orders OrderListResponse
	return &orders, doBLStoreRequest(client, url, "Orders", "orders.json", &orders)
}

func getOrder(client *http.Client, id int64) (*OrderResponse, error) {
	url := fmt.Sprintf(apiBase+"/orders/%d", id)
	var order OrderResponse
	return &order, doBLStoreRequest(client, url, fmt.Sprintf("Order %d", id), fmt.Sprintf("order-%d.json", id), &order)
}

func getColorList(client *http.Client) (*ColorListResponse, error) {
	url := apiBase + "/colors"
	var colors ColorListResponse
	return &colors, doBLStoreRequest(client, url, "Colors", "colors.json", &colors)
}

func getColor(client *http.Client, id int64) (*ColorResponse, error) {
	url := fmt.Sprintf(apiBase+"/colors/%d", id)
	var color ColorResponse
	return &color, doBLStoreRequest(client, url, fmt.Sprintf("Color %d", id), fmt.Sprintf("color-%d.json", id), &color)
}

func doBLStoreRequest(client *http.Client, url, tag, fileName string, v interface{}) error {
	resp, err := client.Get(url)
	printResponseCode(tag, resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeAndWrite(resp.Body, v, fileName)
}

func (r *OrderListResponse) printUnknownValues() {
	for _, o := range r.Orders {
		o.printUnknownValues()
	}
}

func (r *OrderResponse) printUnknownValues() {
	r.Order.printUnknownValues()
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
	Meta   Meta    `json:"meta"`
	Orders []Order `json:"data"`
}

type OrderResponse struct {
	Meta  Meta  `json:"meta"`
	Order Order `json:"data"`
}

type ColorListResponse struct {
	Meta   Meta    `json:"meta"`
	Colors []Color `json:"data"`
}

type ColorResponse struct {
	Meta  Meta  `json:"meta"`
	Color Color `json:"data"`
}

type Order struct {
	OrderID           int64       `json:"order_id"`
	DateOrdered       string      `json:"date_ordered"`
	DateStatusChanged string      `json:"date_status_changed"`
	SellerName        string      `json:"seller_name"`
	StoreName         string      `json:"store_name"`
	BuyerName         string      `json:"buyer_name"`
	BuyerEmail        string      `json:"buyer_email"`       // not in order list
	RequireInsurance  bool        `json:"require_insurance"` // not in order list
	Status            OrderStatus `json:"status"`
	IsInvoiced        bool        `json:"is_invoiced"` // not in order list
	TotalCount        int64       `json:"total_count"`
	UniqueCount       int64       `json:"unique_count"`
	TotalWeight       string      `json:"total_weight"`      // not in order list
	BuyerOrderCount   int64       `json:"buyer_order_count"` // not in order list
	IsFiled           bool        `json:"is_filed"`
	DriveThruSent     bool        `json:"drive_thru_sent"` // not in order list
	Payment           *Payment    `json:"payment"`
	Shipping          *Shipping   `json:"shipping"` // not in order list
	Cost              *Cost       `json:"cost"`
	DisplayCost       *Cost       `json:"disp_cost"`
}

type Payment struct {
	Method       PaymentMethod `json:"method"`
	CurrencyCode CurrencyCode  `json:"currency_code"`
	DatePaid     string        `json:"date_paid,omitempty"`
	Status       PaymentStatus `json:"status"`
}

type Shipping struct {
	MethodID    int64   `json:"method_id"`
	Method      string  `json:"method"`
	Address     Address `json:"address"`
	DateShipped string  `json:"date_shipped"`
}

type Address struct {
	Name        Name   `json:"name"`
	Full        string `json:"full"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
}

type Name struct {
	Full  string `json:"full"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Cost struct {
	CurrencyCode CurrencyCode `json:"currency_code"`
	Subtotal     string       `json:"subtotal"`
	GrandTotal   string       `json:"grand_total"`
	Etc1         string       `json:"etc1,omitempty"`       // not in order list
	Etc2         string       `json:"etc2,omitempty"`       // not in order list
	Insurance    string       `json:"insurance,omitempty"`  // not in order list
	Shipping     string       `json:"shipping,omitempty"`   // not in order list
	Credit       string       `json:"credit,omitempty"`     // not in order list
	Coupon       string       `json:"coupon,omitempty"`     // not in order list
	SalesTax     string       `json:"salesTax,omitempty"`   // not in order list
	VATRate      string       `json:"vat_rate,omitempty"`   // not in order list
	VATAmount    string       `json:"vat_amount,omitempty"` // not in order list
}

// http://apidev.bricklink.com/redmine/projects/bricklink-api/wiki/ColorResource
type Color struct {
	ColorID   int64     `json:"color_id"`   // ID of the color
	ColorName string    `json:"color_name"` // The name of the color
	ColorCode string    `json:"color_code"` // HTML color code of this color
	ColorType ColorType `json:"color_type"` // The name of the color group to which this color belongs
}

type Meta struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Code        int64  `json:"code"`
}

type OrderStatus string
type PaymentStatus string
type PaymentMethod string
type CurrencyCode string
type ColorType string

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
)
