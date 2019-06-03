package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/andrewarchi/brick-apis/bricklinkstore"
	"github.com/andrewarchi/brick-apis/bricklinkuser"
	"github.com/andrewarchi/brick-apis/legobap"
)

var (
	brickLinkStoreConsumerKey    = os.Getenv("BRICKLINK_STORE_CONSUMER_KEY")
	brickLinkStoreConsumerSecret = os.Getenv("BRICKLINK_STORE_CONSUMER_SECRET")
	brickLinkStoreToken          = os.Getenv("BRICKLINK_STORE_TOKEN")
	brickLinkStoreTokenSecret    = os.Getenv("BRICKLINK_STORE_TOKEN_SECRET")
	brickLinkUserUsername        = os.Getenv("BRICKLINK_USER_USERNAME")
	brickLinkUserPassword        = os.Getenv("BRICKLINK_USER_PASSWORD")
	bricksetAPIKey               = os.Getenv("BRICKSET_API_KEY")
	bricksetUsername             = os.Getenv("BRICKSET_USERNAME")
	bricksetPassword             = os.Getenv("BRICKSET_PASSWORD")
	bricksetUserHash             = os.Getenv("BRICKSET_USER_HASH")
	legoBAPAge                   = os.Getenv("LEGO_BAP_AGE")
	legoBAPCountryCode           = os.Getenv("LEGO_BAP_COUNTRY_CODE")
)

func main() {
	blUser, err := bricklinkuser.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	blUser.Login(brickLinkUserUsername, brickLinkUserPassword)
	reportOwnedWantedParts(blUser)

	blStore, err := bricklinkstore.NewClient(brickLinkStoreConsumerKey, brickLinkStoreConsumerSecret, brickLinkStoreToken, brickLinkStoreTokenSecret)
	if err != nil {
		log.Fatal(err)
	}

	orders, err := blStore.GetOrders("out")
	if err != nil {
		fmt.Println(err)
	}
	for _, o := range orders {
		order, err := blStore.GetOrder(o.OrderID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(order)
	}

	items, err := blStore.GetOrderItems(11037590)
	fmt.Println(items)

	age, _ := strconv.Atoi(legoBAPAge)
	bap := legobap.NewClient(age, legobap.CountryCode(legoBAPCountryCode))
	fmt.Println(bap.GetPart("3024"))
	fmt.Println(bap.GetSet("75192"))
}

func printResponse(resp *http.Response, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:", resp.StatusCode, resp.Status)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		str, err := responseToString(resp)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(str)
	}
}

func writeResponse(resp *http.Response, err error, fileName string) {
	fmt.Println("Response:", resp.StatusCode, resp.Status)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	io.Copy(file, resp.Body)
}

func responseToString(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), err
}

func reportOwnedWantedParts(blUser *bricklinkuser.Client) {
	var wanted [][]bricklinkuser.WantedItem
	var sources [][]bricklinkuser.WantedItem
	defaultList, err := blUser.GetWantedList(0)
	if err != nil {
		log.Fatal(err)
	}
	wanted = append(wanted, defaultList.WantedItems)
	for _, list := range defaultList.WantedLists {
		if list.ID != 0 {
			list, err := blUser.GetWantedList(list.ID)
			if err != nil {
				log.Fatal(err)
			}
			name := list.WantedListInfo.Name
			if strings.HasPrefix(name, "[LOOSE]") {
				sources = append(sources, list.WantedItems)
			} else if !strings.HasPrefix(name, "[IGNORE]") {
				wanted = append(wanted, list.WantedItems)
			}
		}
	}
	wantedMap := make(map[string][]bricklinkuser.WantedItem)
	for _, items := range wanted {
		for _, item := range items {
			if item.WantedQty > item.WantedQtyFilled {
				key := item.ItemNumber + ";" + item.ColorName
				wantedMap[key] = append(wantedMap[key], item)
			}
		}
	}
	sourceMap := make(map[string][]bricklinkuser.WantedItem)
	for _, items := range sources {
		for _, item := range items {
			key := item.ItemNumber + ";" + item.ColorName
			sourceMap[key] = append(sourceMap[key], item)
		}
	}
	for _, sourceItems := range sourceMap {
		item := sourceItems[0]
		key := item.ItemNumber + ";" + item.ColorName
		if wantedItems, ok := wantedMap[key]; ok {
			fmt.Println("Want", item.ItemType, item.ColorName, item.ItemNumber, item.ItemName)
			for _, wi := range wantedItems {
				fmt.Printf(" - %d (have %d) in %s\n", wi.WantedQty-wi.WantedQtyFilled, wi.WantedQtyFilled, wi.WantedListName)
			}
			fmt.Println("Sources")
			for _, si := range sourceItems {
				name := strings.TrimSpace(strings.TrimPrefix(si.WantedListName, "[LOOSE]"))
				fmt.Printf(" - %d in %s\n", si.WantedQty, name)
			}
			fmt.Println()
		}
	}
}
