package brickset

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/fiorix/wsdl2go/soap"
)

const (
	apiKey   = ""
	url      = "https://brickset.com/api"
	endpoint = url + "/v2.asmx"
	username = ""
	password = ""
)

type soapClient struct {
	c *http.Client
}

func getFormEncodedEndpoint(methodName string) string {
	return fmt.Sprintf("%s/%s", endpoint, methodName)
}

func getSOAPAction(methodName string) string {
	return fmt.Sprintf("%s/%s", url, methodName)
}

func getXML(apiKey, username, password string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<login xmlns="https://brickset.com/api/">
		  <apiKey>%s</apiKey>
		  <username>%s</username>
		  <password>%s</password>
		</login>
	  </soap:Body>
	</soap:Envelope>`, apiKey, username, password)
}

func TestLogin(t *testing.T) {
	svc := NewAPIv2Service(&soap.Client{
		URL:       url,
		Namespace: Namespace,
	})
	t.Fatal(svc.Login(apiKey, username, password))
}

func TestLogin2(t *testing.T) {
	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(getXML(apiKey, username, password)))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", getSOAPAction("login"))
	fmt.Println(http.DefaultClient.Do(req))
}

func TestLogin3(t *testing.T) {
	c := newClient()
	t.Fatal(c.login(apiKey, username, password))
}

func newClient() *soapClient {
	return &soapClient{&http.Client{}}
}

func (c soapClient) login(apiKey, username, password string) (string, error) {
	body := strings.NewReader(fmt.Sprintf(`apiKey=%s&username=%s&password=%s`, apiKey, username, password))
	req, _ := http.NewRequest("POST", getFormEncodedEndpoint("login"), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	resp, err := c.c.Do(req)
	if err != nil {
		return "", err
	}
	return decodeLoginResponse(resp)
}

func decodeLoginResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	r := &loginResponse{}
	if err := xml.NewDecoder(resp.Body).Decode(r); err != nil {
		return "", err
	}
	return r.Response, nil
}

type loginResponse struct {
	XMLName  xml.Name `xml:"https://brickset.com/api/ string"`
	Response string   `xml:",chardata"`
}

func TestDecodeLoginResponseXML(t *testing.T) {
	xmlString := `<string xmlns="https://brickset.com/api/">test</string>`
	r := &loginResponse{}
	if err := xml.NewDecoder(strings.NewReader(xmlString)).Decode(r); err != nil {
		t.Fatal(err)
	}
	if r.Response != "test" {
		t.Error("Expected correct string value", r.Response)
	}
}

func TestEncodeXml(t *testing.T) {
	data, err := xml.Marshal(&loginResponse{Response: "value"})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `<string xmlns="https://brickset.com/api/">test</string>` {
		t.Error("unexpected result")
	}
}
