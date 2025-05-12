package api

type LinksResponse struct {
	Links []Link `json:"links"`
}

type Link struct {
	URL      string `json:"irl"`
	Category string `json:"category"`
}
