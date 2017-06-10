package OnurTPIItem

//Item structure
type Item struct {
	ID         string   `json:"id,omitempty"`
	Title      string   `json:"title,omitempty"`
	Price      float64  `json:"price"`
	ImgURL     string   `json:"imUrl,omitempty"`
	Related    []string `json:"related,omitempty"`
	SalesRank  string   `json:"salesRank,omitempty"`
	Brand      string   `json:"brand,omitempty"`
	Categories []string `json:"categories,omitempty"`
}
