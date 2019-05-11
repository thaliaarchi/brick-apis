package main

import (
	"fmt"
	"log"
	"net/http"
)

func getBricksAndPiecesProduct(cred *credentials, part string) (*http.Response, error) {
	url := "https://www.lego.com/en-US/service/rpservice/getproduct?productnumber=" + part + "&isSalesFlow=true"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	cookie := fmt.Sprintf(`csAgeAndCountry={"age":"%s","countrycode":"%s"}`, cred.Age, cred.CountryCode)
	request.Header.Add("Cookie", cookie)
	return http.DefaultClient.Do(request)
}
