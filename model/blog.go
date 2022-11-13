package model

type Blog struct {
	Code      string     `json:"code"`
	Title     string     `json:"title"`
	Text      string     `json:"text"`
	Img       string     `json:"img"`
	Date      string     `json:"date"`
	Link      string     `json:"link"`
	ArtiCode  string     `json:"arti_code"`
	ArtistImg string     `json:"artist_img"`
	Name      string     `json:"name"`
	Comments  []*Comment `json:"comments"`
}

type Blogs struct {
	Count string  `json:"count"`
	Data  []*Blog `json:"data"`
}
