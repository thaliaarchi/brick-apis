package bricklinkuser

import (
	"testing"

	"github.com/andrewarchi/brick-apis/credentials"
)

func TestRetrieveCartInfo(t *testing.T) {
	cred, err := credentials.Read("../credentials.json")
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(cred.BrickLinkUser)
	client.Login()
	_, err = client.RetrieveCartInfo()
	if err != nil {
		t.Fatal(err)
	}
}
