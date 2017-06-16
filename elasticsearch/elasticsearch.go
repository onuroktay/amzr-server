package OnurTPIES

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"net/http"
	"strings"
)

const (
	ITEM = "item"
	USER = "user"
)

type ElasticSearch struct {
	client     *elastic.Client
	_indexName string
}

var testMode bool

func SetTestMode(status bool) {
	testMode = status
}

func getESType(defaultType string) string {
	if testMode {
		return "test"
	}

	return defaultType
}

// NewElasticSearch open a connection to the ElasticSearch database
func NewElasticSearch(indexName string) (es *ElasticSearch, err error) {
	if indexName == "" {
		err = errors.New("Index Name is mandatory")
		return
	}

	es = new(ElasticSearch)

	es._indexName = indexName
	es.client, err = elastic.NewClient()
	if err != nil {
		return
	}

	err = es.createIndexIfNotExist()

	return
}

func (es *ElasticSearch) createIndexIfNotExist() error {
	exists, err := es.client.IndexExists(es._indexName).Do(context.TODO())
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	createIndex, err := es.client.CreateIndex(es._indexName).Do(context.TODO())
	if err != nil {
		return err
	}

	if !createIndex.Acknowledged {
		err = errors.New(fmt.Sprintf("expected IndicesCreateResult.Acknowledged %v; got %v", true, createIndex.Acknowledged))
	}

	return err
}

func (es *ElasticSearch) executeQuery(_type, search, query string) ([]byte, error) {

	fmt.Println("query 2b exec:",query)
	// Transform string into io.Reader
	body := strings.NewReader(query)

	//fmt.Println("http://localhost:9200/amazonreader/account/_search?" + search)
	req, err := http.NewRequest("GET", "http://localhost:9200/amazonreader/"+_type+"/_search?"+search, body)

	if err != nil {
		// handle err
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		return nil, err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}
