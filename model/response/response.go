package response

type SearchResponse struct {
	Hits Hits `json:"hits,omitempty"`
}

type Hits struct {
	Total    Total   `json:"total,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Hits     []Hit   `json:"hits,omitempty"`
}

type Total struct {
	Value    int    `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type Hit struct {
	Score     float64             `json:"_score,omitempty"`
	Source    Source              `json:"_source,omitempty"`
	Highlight map[string][]string `json:"highlight,omitempty"`
}

type Source struct {
	Title string `json:"title,omitempty"`
	Name  string `json:"name,omitempty"`
	Date  string `json:"date,omitempty"`
	Img   string `json:"img,omitempty"`
	Link  string `json:"link,omitempty"`
}
