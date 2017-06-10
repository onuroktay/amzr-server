package OnurTPIES

type Search struct {
	Category  string   `json:"category"`
	Title     string   `json:"title"`
	PriceFrom string   `json:"priceFrom"`
	PriceTo   string   `json:"priceTo"`
	Size      int   `json:"size"`
	From      int   `json:"from"`
}
