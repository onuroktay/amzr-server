package OnurTPIItem

// SubHits contains a sub-data structure returns by elasticsearch
type SubHits struct {
	Source Item   `json:"_source"`
	ID     string `json:"_id"`
}

type MainHits struct {
	Hits  []SubHits `json:"hits"`
	Total int       `json:"total"`
}

type ESResponse struct {
	Took    int      `json:"took"`
	TimeOut bool     `json:"time_out"`
	Hits    MainHits `json:"hits"`
}

type SearchResponse struct {
	Items []*Item `json:"items"`
	Total int       `json:"total"`
}
