package convertor

import "server-scrapper/internal/lib/api"

func ConvertLinks(userLinks map[string]string) []api.Link {
	var links []api.Link
	for url, category := range userLinks {
		links = append(links, api.Link{
			URL:      url,
			Category: category,
		})
	}
	return links
}
