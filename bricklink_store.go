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

func getOrderDetails(client *http.Client, id int64) (*http.Response, error) {
	url := fmt.Sprintf(apiBase+"/orders/%d", id)
	return client.Get(url)
}
