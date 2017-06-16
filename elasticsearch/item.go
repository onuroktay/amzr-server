package OnurTPIES

import (
	"encoding/json"
	"errors"
	"github.com/onuroktay/amazon-reader/amzr-server/item-data"
	"github.com/onuroktay/amazon-reader/amzr-server/util"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func (es *ElasticSearch) GetItemByIDInDB(id string) (*OnurTPIItem.Item, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}

	var query = `
{
  "query": {
    "term" : { "_id" :"` + id + `"  }
  }
}`

	res, err := es.executeQuery(getESType(ITEM), "", query)
	if err != nil {
		return nil, err
	}

	response := &OnurTPIItem.ESResponse{}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	if response != nil && len(response.Hits.Hits) > 0 {
		item := response.Hits.Hits[0].Source

		return &item, nil
	}

	return nil, nil
}

// SaveItem save an item in ElasticSearch
func (es *ElasticSearch) SaveItem(item *OnurTPIItem.Item) (err error) {
	_, err = es.client.Index().
		Index(es._indexName).
		Type(getESType(ITEM)).
		Id(item.ID).
		BodyJson(item).
		Do(context.TODO())

	return
}

func (es *ElasticSearch) UpdateItem(id string, item *OnurTPIItem.Item) error {
	oldItem, err := es.GetItemByIDInDB(id)
	if err != nil {
		// Handle error
		return err
	}

	// update item
	oldItem.Title = item.Title
	oldItem.Price = item.Price
	oldItem.ImgURL = item.ImgURL

	// save updated item
	_, err = es.client.Index().
		Index(es._indexName).
		Type(getESType(ITEM)).
		Id(id).
		BodyJson(oldItem).
		Do(context.Background())

	if err != nil {
		// Handle error
		log.Println(err)
		return err
	}

	return err
}

func (es *ElasticSearch) DeleteItem(id string) error {
	_, err := es.client.Delete().
		Index(es._indexName).
		Type(getESType(ITEM)).
		Id(id).
		Refresh("true").
		Do(context.TODO())

	return err
}

func (es *ElasticSearch) GetItems(criteria *Search) (*OnurTPIItem.SearchResponse, error) {

	// Default query : all items
	var query = mainQuery(criteria)

	res, err := es.executeQuery(getESType(ITEM), "", query)
	if err != nil {
		return nil, err
	}

	response := &OnurTPIItem.ESResponse{}
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	if response != nil {
		var items []*OnurTPIItem.Item

		for _, hit := range response.Hits.Hits {
			item := hit.Source
			item.ID = hit.ID
			items = append(items, &item)
		}

		resp := &OnurTPIItem.SearchResponse{
			Items: items,
			Total: response.Hits.Total,
		}

		return resp, nil
	}

	return nil, nil

}

func queryMust(criteria *Search) (queryMust string) {
	queryCategory := queryCategory(criteria)
	queryTitle := queryTitle(criteria)

	query := queryCategory
	if queryCategory != "" && queryTitle != "" {
		query += ","
	}
	query += queryTitle

	if query != "" {
		queryMust = `
		"must": [` +
			query +
			`]`
	}

	return
}

func queryCategory(criteria *Search) string {
	if criteria.Category == "" {
		return ""
	}

	return `
	{
          "match_phrase": {
            "categories": {
              "query": "` + util.CleanQuote(criteria.Category) + `"
            }
          }
        }
        `
}

func queryTitle(criteria *Search) string {
	if criteria.Title == "" {
		return ""
	}

	return `
	{
          "match_phrase": {
            "title": {
              "query": "` + util.CleanQuote(criteria.Title) + `"
            }
          }
        }
        `
}

func queryFilter(criteria *Search) string {
	queryPriceFrom := filterPriceFrom(criteria)
	queryPriceTo := filterPriceTo(criteria)

	query := queryPriceFrom
	if queryPriceFrom != "" && queryPriceTo != "" {
		query += ","
	}
	query += queryPriceTo

	if query == "" {
		return ""
	}

	return `
	"filter": {
	  "range": {
	    "price": {` +
		query + `
	    }
	  }
	}
	`
}

func filterPriceFrom(criteria *Search) string {
	priceFrom, _ := strconv.ParseFloat(criteria.PriceFrom, 64)
	if priceFrom <= 0 {
		return ""
	}

	return `"gte": ` + criteria.PriceFrom
}

func filterPriceTo(criteria *Search) string {
	priceFrom, _ := strconv.ParseFloat(criteria.PriceTo, 64)
	if priceFrom <= 0 {
		return ""
	}

	return `"lte": ` + criteria.PriceTo
}

func queryBool(criteria *Search) string {
	queryMust := queryMust(criteria)
	queryFilter := queryFilter(criteria)

	query := queryMust
	if queryMust != "" && queryFilter != "" {
		query += ","
	}
	query += queryFilter

	return `
	"bool": {` +
		query +
		`
		}
		`
}

func querySort(criteria *Search) string {
	return `
	"sort": [
	   {
	     "price": {
	     	"order": "asc"
	     }
	   }
	],
	"size": ` + strconv.Itoa(criteria.Size) + `,
	"from": ` + strconv.Itoa(criteria.From) + `
	`
}

func mainQuery(criteria *Search) (query string) {
	queryBool := queryBool(criteria)
	querySort := querySort(criteria)

	if queryBool != "" {
		query = `
		  "query": {` +
			queryBool + `
		  },
		  `
	}
	query += querySort

	return `
	{` +
		query + `
	}`
}

