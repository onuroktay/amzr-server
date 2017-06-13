package OnurTPIES

import (
	"encoding/json"
	"errors"
	"github.com/onuroktay/amazon-reader/amzr-server/user-data"
	"github.com/onuroktay/amazon-reader/amzr-server/util"
	"golang.org/x/net/context"
)

func (es *ElasticSearch) DeleteAccount(id string) error {
	var nbAdmin int
	var isAdmin bool


	// Get users
	users, err := es.GetUsers()

	if err != nil {
		return err
	}

	// Look for Admin
	for _, user := range users {
		// Check if last Admin
		if user.RoleValue == OnurTPIUser.ADMIN {
			nbAdmin++
		}

		// Check if user to delete is Admin
		if user.ID == id && user.RoleValue == OnurTPIUser.ADMIN  {
			isAdmin = true
		}
	}

	// Check to have at least one admin
	if isAdmin && nbAdmin == 1{
		return errors.New("Sorry, last admin can't be deleted")
	}

	_, err = es.client.Delete().
		Index(es._indexName).
		Type(getESType(USER)).
		Id(id).
		Refresh("true").
		Do(context.TODO())

	return err
}

// SaveAccount save an acount in ElasticSearch
func (es *ElasticSearch) CreateAccount(account *OnurTPIUser.Account) (string, error) {
	put1, err := es.client.Index().
		Index(es._indexName).
		Type(getESType(USER)).
		Id(util.GetUUID()).
		BodyJson(account).
		Do(context.Background())

	if err != nil {
		// Handle error
		return "", err
	}

	// Information available = put1.Id, put1.Index, put1.Type

	return put1.Id, nil
}

func (es *ElasticSearch) UpdateRole(id string, roleValue int) error {
	account, err := es.GetAccountByIDInDB(id)
	if err != nil {
		// Handle error
		return err
	}

	user := OnurTPIUser.Account{}
	user.UserName = account.UserName
	user.RoleValue = roleValue
	user.Password = account.Password

	_, err = es.client.Index().
		Index(es._indexName).
		Type(getESType(USER)).
		Id(id).
		BodyJson(user).
		Do(context.Background())

	if err != nil {
		// Handle error
		return err
	}

	return err

}

func (es *ElasticSearch) CheckIfUserExistsInDB(cred *OnurTPIUser.CredentialsClient) (bool, error) {
	if cred == nil || cred.UserName == "" {
		return false, errors.New("missing username")
	}

	source, err := es.GetAccountByUserNameInDB(cred)
	if err != nil {
		return false, err
	}

	if source != nil && source.ID != "" {
		return true, nil
	}

	return false, nil
}

func (es *ElasticSearch) GetAccountByUserNameInDB(cred *OnurTPIUser.CredentialsClient) (*OnurTPIUser.Account, error) {
	if cred == nil || util.CleanQuote(cred.UserName) == "" {
		return nil, errors.New("missing username")
	}

	var query = `
{
  "query": {
    "match_phrase" : { "username" :"` + util.CleanQuote(cred.UserName) + `"  }
  }
}`

	res, err := es.executeQuery(getESType(USER), "", query)
	if err != nil {
		return nil, err
	}

	response := &OnurTPIUser.ESResponse{}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	if response != nil && len(response.Hits.Hits) > 0 {
		account := response.Hits.Hits[0].Source
		account.ID = response.Hits.Hits[0].ID

		return &account, nil
	}

	return nil, nil
}

func (es *ElasticSearch) GetAccountByIDInDB(id string) (*OnurTPIUser.Account, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}

	var query = `
{
  "query": {
    "match_phrase" : { "_id" :"` + id + `"  }
  }
}`

	res, err := es.executeQuery(getESType(USER), "", query)
	if err != nil {
		return nil, err
	}

	response := &OnurTPIUser.ESResponse{}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	if response != nil && len(response.Hits.Hits) > 0 {
		account := response.Hits.Hits[0].Source
		account.ID = response.Hits.Hits[0].ID

		return &account, nil
	}

	return nil, errors.New("account not found")
}

func (es *ElasticSearch) GetUsers() ([]*OnurTPIUser.User, error) {
	var query = `
{
}`

	res, err := es.executeQuery(getESType(USER), "", query)
	if err != nil {
		return nil, err
	}

	response := &OnurTPIUser.ESResponse{}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	if response != nil && len(response.Hits.Hits) > 0 {
		var users []*OnurTPIUser.User

		for _, hit := range response.Hits.Hits {
			user := hit.Source.User
			user.ID = hit.ID
			users = append(users, &user)
		}

		return users, nil
	}

	return nil, nil
}
