package main

import (
	"testing"
	"log"
	"github.com/onuroktay/amazon-reader/AmzR-Server/elasticsearch"
	"github.com/onuroktay/amazon-reader/AmzR-Server/user-data"
	"github.com/onuroktay/amazon-reader/AmzR-Server/item-data"
)

func init() {
	// Connect to ElasticSearch
	es, err := OnurTPIES.NewElasticSearch("amazonreader")

	if err != nil {
		log.Fatalln("ElasticSearch connection error:", err.Error())
	}

	// Set Database
	database = &DATABASE{accesser: es}

	// Set Test Mode
	OnurTPIES.SetTestMode(true)
}

func TestUserFunctions(t *testing.T) {

	// Create account
	account := &OnurTPIUser.Account{}
	account.UserName = "test1"
	account.RoleValue = 2
	account.Password = "test"

	database.accesser.CreateAccount(account)

	// Read account in DB
	cred := &OnurTPIUser.CredentialsClient{}
	cred.UserName = "test1"
	cred.Password = "test"

	readAccount, err := database.accesser.GetAccountByUserNameInDB(cred)
	if err != nil {
		t.Error()
	}

	// Check account
	if readAccount == nil {
		t.Error()
	}

	if readAccount.UserName != account.UserName ||
		readAccount.Password != account.Password {
		t.Error()
	}

	// SaveUser sets rolevalue to 0
	if readAccount.RoleValue != 0 {
		t.Error()
	}

	// Delete Account
	err = database.accesser.DeleteAccount(readAccount.ID)
	if err != nil {
		t.Error()
	}
}

func TestItemsFunctions(t *testing.T) {

	// Create item
	item := &OnurTPIItem.Item{
		ID:         "item1",
		Title:      "item",
		Price:      100,
		ImgURL:     "",
		Categories: []string{"test"},
	}

	database.accesser.SaveItem(item)

	// Read account in DB
	cred := &OnurTPIItem.Item{}
	cred.ID = "item1"

	readItem, err := database.accesser.GetItemByIDInDB(cred.ID)
	if err != nil {
		t.Error()
	}

	// Check item
	if readItem == nil {
		t.Error()
	}

	if readItem.ID != item.ID ||
		readItem.Title != item.Title {
		t.Error()
	}

	// Delete Account
	err = database.accesser.DeleteItem(readItem.ID)
	if err != nil {
		t.Error()
	}

}
