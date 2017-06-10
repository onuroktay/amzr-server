package main

import (
	"github.com/onuroktay/amazon-reader/AmzR-Server/elasticsearch"
	"github.com/onuroktay/amazon-reader/AmzR-Server/item-data"
	"github.com/onuroktay/amazon-reader/AmzR-Server/user-data"
)

// DATABASE contains an interface with the method to implement in order to use de database
type DATABASE struct {
	// List of Methods to be implemented in the db struct (couchbase, elasticsearch, ...)
	accesser interface {
		DeleteAccount(id string) error
		CreateAccount(account *OnurTPIUser.Account) (string, error)
		UpdateRole(id string, roleValue int) error
		CheckIfUserExistsInDB(cred *OnurTPIUser.CredentialsClient) (bool, error)
		GetAccountByUserNameInDB(cred *OnurTPIUser.CredentialsClient) (*OnurTPIUser.Account, error)
		GetAccountByIDInDB(id string) (*OnurTPIUser.Account, error)
		GetUsers() ([]*OnurTPIUser.User, error)

		SaveItem(item *OnurTPIItem.Item) (err error)
		GetItemByIDInDB(id string) (*OnurTPIItem.Item, error)
		UpdateItem(id string, item *OnurTPIItem.Item) error
		DeleteItem(id string) error
		GetItems(searchFromClient *OnurTPIES.Search) (*OnurTPIItem.SearchResponse, error)
	}
}
