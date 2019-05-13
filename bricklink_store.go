package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/mrjones/oauth"
)

const (
	apiBase = "https://api.bricklink.com/api/store/v1"
)

type BrickLinkStoreClient struct {
	client *http.Client
}

func NewBrickLinkStoreClient(cred *BrickLinkCredentials) (*BrickLinkStoreClient, error) {
	consumer := oauth.NewConsumer(cred.ConsumerKey, cred.ConsumerSecret, oauth.ServiceProvider{})
	accessToken := &oauth.AccessToken{Token: cred.Token, Secret: cred.TokenSecret}
	client, err := consumer.MakeHttpClient(accessToken)
	return &BrickLinkStoreClient{client}, err
}

func (c *BrickLinkStoreClient) GetOrderList() ([]Order, error) {
	url := apiBase + "/orders?direction=out"
	var orders OrderListResponse
	if err := c.doRequest(url, "Orders", "orders.json", &orders); err != nil {
		return nil, err
	}
	return orders.Data, nil
}

func (c *BrickLinkStoreClient) GetOrder(id int64) (*Order, error) {
	url := fmt.Sprintf(apiBase+"/orders/%d", id)
	var order OrderResponse
	if err := c.doRequest(url, fmt.Sprintf("Order %d", id), fmt.Sprintf("order-%d.json", id), &order); err != nil {
		return nil, err
	}
	return &order.Data, nil
}

func (c *BrickLinkStoreClient) GetColorList() ([]Color, error) {
	url := apiBase + "/colors"
	var colors ColorListResponse
	if err := c.doRequest(url, "Colors", "colors.json", &colors); err != nil {
		return nil, err
	}
	return colors.Data, nil
}

func (c *BrickLinkStoreClient) GetColor(id int64) (*Color, error) {
	url := fmt.Sprintf(apiBase+"/colors/%d", id)
	var color ColorResponse
	if err := c.doRequest(url, fmt.Sprintf("Color %d", id), fmt.Sprintf("color-%d.json", id), &color); err != nil {
		return nil, err
	}
	return &color.Data, nil
}

func (c *BrickLinkStoreClient) doRequest(url, tag, fileName string, v interface{}) error {
	resp, err := c.client.Get(url)
	printResponseCode(tag, resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := decodeAndWrite(resp.Body, v, fileName); err != nil {
		return err
	}
	f := reflect.ValueOf(v).Elem().FieldByName("Meta")
	if f != reflect.ValueOf(nil) {
		meta := f.Interface().(Meta)
		if meta.Description != "OK" || meta.Message != "OK" || meta.Code != 200 {
			return fmt.Errorf("Meta is not OK 200: %v", meta)
		}
	}
	return nil
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

type ColorListResponse struct {
	Meta Meta    `json:"meta"`
	Data []Color `json:"data"`
}

type ColorResponse struct {
	Meta Meta  `json:"meta"`
	Data Color `json:"data"`
}

type Order struct {
	OrderID           int64       `json:"order_id"`
	DateOrdered       time.Time   `json:"date_ordered"`
	DateStatusChanged time.Time   `json:"date_status_changed"`
	SellerName        string      `json:"seller_name"`
	StoreName         string      `json:"store_name"`
	BuyerName         string      `json:"buyer_name"`
	BuyerEmail        string      `json:"buyer_email"`       // not in order list
	RequireInsurance  bool        `json:"require_insurance"` // not in order list
	Status            OrderStatus `json:"status"`
	IsInvoiced        bool        `json:"is_invoiced"` // not in order list
	TotalCount        int64       `json:"total_count"`
	UniqueCount       int64       `json:"unique_count"`
	TotalWeight       float64     `json:"total_weight,string"` // not in order list
	BuyerOrderCount   int64       `json:"buyer_order_count"`   // not in order list
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
	DatePaid     time.Time     `json:"date_paid,omitempty"`
	Status       PaymentStatus `json:"status"`
}

type Shipping struct {
	MethodID    int64     `json:"method_id"`
	Method      string    `json:"method"`
	Address     Address   `json:"address"`
	DateShipped time.Time `json:"date_shipped"`
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
	Subtotal     float64      `json:"subtotal,string"`
	GrandTotal   float64      `json:"grand_total,string"`
	Etc1         float64      `json:"etc1,string"`       // not in order list
	Etc2         float64      `json:"etc2,string"`       // not in order list
	Insurance    float64      `json:"insurance,string"`  // not in order list
	Shipping     float64      `json:"shipping,string"`   // not in order list
	Credit       float64      `json:"credit,string"`     // not in order list
	Coupon       float64      `json:"coupon,string"`     // not in order list
	SalesTax     float64      `json:"salesTax,string"`   // not in order list
	VATRate      float64      `json:"vat_rate,string"`   // not in order list
	VATAmount    float64      `json:"vat_amount,string"` // not in order list
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
