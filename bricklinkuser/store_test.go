package bricklinkuser

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	username = os.Getenv("BRICKLINK_USER_USERNAME")
	password = os.Getenv("BRICKLINK_USER_PASSWORD")
)

func TestAddToCart(t *testing.T) {
	t.SkipNow()
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Login(username, password); err != nil {
		t.Fatal(err)
	}
	t.Error(c.AddToCart("596847", []CartItemSimple{{ID: 170802686, Quantity: "1", SellerID: 596847, SourceType: 1}}))
}

func TestGetAddToCartQuery(t *testing.T) {
	query, err := getAddToCartQuery("596847", []CartItemSimple{{ID: 170802686, Quantity: "1", SellerID: 596847, SourceType: 1}})
	if err != nil {
		t.Fatal(err)
	}
	expected := url.Values{}
	expected.Add("sid", "596847")
	expected.Add("itemArray", `[{"invID":170802686,"invQty":"1","sellerID":596847,"sourceType":1}]`)
	if !reflect.DeepEqual(expected, query) {
		t.Errorf("Values don't match Expected:%v, Actual:%v", expected, query)
	}
}

func TestConvertCartReturnJSON(t *testing.T) {
	jsonText := `{"errors":0,"itemReturnStatus":[{"invID":170057995,"code":"0","msg":"OK","sid":1189138}],"carts":[{"sellerID":1189138,"vatPct":1,"sellerName":"djinn","storeName":"Me and my bricks","countryID":"NO","feedback":22,"current_cart":{"items":[{"itemName":"Death Star","invDescription":"Deathstar has all bricks accounted for, including building-instructions.\nWithout minifigures or weapons.","invID":170057995,"invQty":1,"bulkQty":1,"superlotID":0,"superlotQty":1,"salePercent":0,"itemType":"S","itemBrand":1000,"invNew":"Used","invComplete":"Incomplete","colorID":0,"colorName":"(Not Applicable)","itemNo":"10188","itemSeq":1,"itemID":78848,"itemStatus":"A","smallImg":"//img.bricklink.com/ItemImage/ST/0/10188-1.t1.png","largeImg":"//img.bricklink.com/ItemImage/SN/0/10188-1.png","nativePrice":"NOK 1,900.00","salePrice":"US $218.7587","invPrice":"US $218.7587","invTierQty1":0,"invTierPrice1":"US $0.00","invTierSalePrice1":"US $0.00","invTierNativeSalePrice1":"NOK 0.00","invTierQty2":0,"invTierPrice2":"US $0.00","invTierSalePrice2":"US $0.00","invTierNativeSalePrice2":"NOK 0.00","invTierQty3":0,"invTierPrice3":"US $0.00","invTierSalePrice3":"US $0.00","invTierNativeSalePrice3":"NOK 0.00","cartQty":1,"cartBindQty":1,"invDate":"2019-05-25","invASCAvailable":true,"invAvailable":"Y","warnings":[],"totalWeightOunces":284.31,"totalWeightGrams":"8060.0","totalPrice":"US $218.76","totalSalePrice":"US $218.76","totalNativePrice":"NOK 1,900.00","totalNativeSalePrice":"NOK 1,900.00"}],"superlots":[],"totalItems":1,"totalLots":1,"totalPrice":"US $218.76","totalNativePrice":"NOK 1,900.00","totalWarnings":0,"totalNativePriceRaw":"1900.0000","totalWeightGrams":"8060.0","totalWeightOunces":"284.31","weightUnknownLots":0,"aveLotPrice":"NOK 1,900.00"}}],"totStoreCartCnt":2,"cartItemErrorCode":0,"returnCode":0,"returnMessage":"OK","errorTicket":0,"procssingTime":38}`
	r, err := decodeCartReturn(ioutil.NopCloser(strings.NewReader(jsonText)))
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(r)
	compareStrings(string(data), jsonText, t)

	expected := AddToCartResponse{
		Errors:           0,
		ItemReturnStatus: []ItemReturnStatus{{170057995, "0", "OK", 1189138}},
		Carts: []StoreCart{{
			SellerID:   1189138,
			VATPct:     1,
			SellerName: "djinn",
			StoreName:  "Me and my bricks",
			CountryID:  "NO",
			Feedback:   22,
			CurrentCart: CartItems{
				Items: []CartItemDetail{{
					ItemName:                      "Death Star",
					InventoryDescription:          "Deathstar has all bricks accounted for, including building-instructions.\nWithout minifigures or weapons.",
					InventoryID:                   170057995,
					InventoryQuantity:             1,
					BulkQuantity:                  1,
					SuperLotID:                    0,
					SuperLotQuantity:              1,
					SalePercent:                   0,
					ItemType:                      "S",
					ItemBrand:                     1000,
					InventoryCondition:            "Used",
					InventoryComplete:             "Incomplete",
					ColorID:                       0,
					ColorName:                     "(Not Applicable)",
					ItemNumber:                    "10188",
					ItemSequence:                  1,
					ItemID:                        78848,
					ItemStatus:                    "A",
					SmallImage:                    "//img.bricklink.com/ItemImage/ST/0/10188-1.t1.png",
					LargeImage:                    "//img.bricklink.com/ItemImage/SN/0/10188-1.png",
					NativePrice:                   "NOK 1,900.00",
					SalePrice:                     "US $218.7587",
					InventoryPrice:                "US $218.7587",
					InventoryTierQuantity1:        0,
					InventoryTierPrice1:           "US $0.00",
					InventoryTierSalePrice1:       "US $0.00",
					InventoryTierNativeSalePrice1: "NOK 0.00",
					InventoryTierQuantity2:        0,
					InventoryTierPrice2:           "US $0.00",
					InventoryTierSalePrice2:       "US $0.00",
					InventoryTierNativeSalePrice2: "NOK 0.00",
					InventoryTierQuantity3:        0,
					InventoryTierPrice3:           "US $0.00",
					InventoryTierSalePrice3:       "US $0.00",
					InventoryTierNativeSalePrice3: "NOK 0.00",
					CartQuantity:                  1,
					CartBindQuantity:              1,
					InventoryDate:                 "2019-05-25",
					InventoryASCAvailable:         true,
					InventoryAvailable:            "Y",
					Warnings:                      nil,
					TotalWeightOunces:             284.31,
					TotalWeightGrams:              "8060.0",
					TotalPrice:                    "US $218.76",
					TotalSalePrice:                "US $218.76",
					TotalNativePrice:              "NOK 1,900.00",
					TotalNativeSalePrice:          "NOK 1,900.00",
				}},
				Superlots:           nil,
				TotalItems:          1,
				TotalLots:           1,
				TotalPrice:          "US $218.76",
				TotalNativePrice:    "NOK 1,900.00",
				TotalWarnings:       0,
				TotalNativePriceRaw: "1900.0000",
				TotalWeightGrams:    "8060.0",
				TotalWeightOunces:   "284.31",
				WeightUnknownLots:   0,
				AverageLotPrice:     "NOK 1,900.00",
			},
		}},
		TotalStoreCartCount: 2,
		CartItemErrorCode:   0,
		ReturnCode:          0,
		ReturnMessage:       "OK",
		ErrorTicket:         0,
		ProcessingTime:      38,
	}

	if reflect.DeepEqual(r, expected) {
		t.Error("expected match\n", r, "\n", expected)
	}
}

func compareStrings(s1, s2 string, t *testing.T) {
	if len(s1) != len(s2) {
		t.Error("string length doesn't match", len(s1), len(s2))
	}
	for i := 0; i < len(s1) && i < len(s2); i++ {
		if s1[i] != s2[i] {
			start := i - 10
			if start < 0 {
				start = 0
			}
			t.Error("mismatch found starting at character", i)
			t.Error("s1:", s1[start:i+1])
			t.Error("s2:", s2[start:i+1])
			return
		}
	}
}
