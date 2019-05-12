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

func getOrders(client *http.Client) (*OrderListResponse, error) {
	resp, err := client.Get(apiBase + "/orders?direction=out")
	printResponseCode("Orders", resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var orders OrderListResponse
	err = decodeAndWrite(resp.Body, &orders, "orders.json")
	return &orders, err
}

func getOrderDetails(client *http.Client, id int64) (*OrderResponse, error) {
	resp, err := client.Get(fmt.Sprintf(apiBase+"/orders/%d", id))
	printResponseCode(fmt.Sprintf("Order %d", id), resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var order OrderResponse
	err = decodeAndWrite(resp.Body, &order, fmt.Sprintf("order-%d.json", id))
	return &order, err
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
		case PaymentPayPal:
		case PaymentPayPalOnsite:
		}
		switch o.Payment.CurrencyCode {
		default:
			fmt.Printf("CurrencyCode: %s", o.Payment.CurrencyCode)
		case CurrencyCad:
		case CurrencyEur:
		case CurrencyHuf:
		case CurrencyUsd:
		}
		switch o.Payment.Status {
		default:
			fmt.Printf("PaymentStatus: %s", o.Payment.Status)
		case PaymentCompleted:
		case PaymentReceived:
		case PaymentNone:
		case PaymentSent:
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

type Meta struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Code        int64  `json:"code"`
}

type OrderStatus string
type PaymentStatus string
type PaymentMethod string
type CurrencyCode string

const (
	OrderCompleted      OrderStatus   = "COMPLETED"
	OrderReceived       OrderStatus   = "RECEIVED"
	OrderPurged         OrderStatus   = "PURGED"
	PaymentCompleted    PaymentStatus = "Completed"
	PaymentReceived     PaymentStatus = "Received"
	PaymentNone         PaymentStatus = "None"
	PaymentSent         PaymentStatus = "Sent"
	PaymentPayPal       PaymentMethod = "PayPal"
	PaymentPayPalOnsite PaymentMethod = "PayPal (Onsite)"
	CurrencyCad         CurrencyCode  = "CAD"
	CurrencyEur         CurrencyCode  = "EUR"
	CurrencyHuf         CurrencyCode  = "HUF"
	CurrencyUsd         CurrencyCode  = "USD"
)
