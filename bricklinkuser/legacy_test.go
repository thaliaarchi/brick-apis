package bricklinkuser

import "testing"

func TestGetCartInfo(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	err = client.Login(username, password)
	if err != nil {
		t.Fatal(err)
	}
	cart, err := client.GetGlobalCart()
	if err != nil {
		t.Fatal(err)
	}
	t.Error(cart)
}

func TestGetCheckoutInfo(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	err = client.Login(username, password)
	if err != nil {
		t.Fatal(err)
	}
	checkout, err := client.GetGlobalCartCheckoutInfo(812889, "812889:-536904708:1551047564403")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(checkout)
}
