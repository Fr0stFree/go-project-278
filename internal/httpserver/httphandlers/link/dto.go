package link

import (
	"shortener/internal/services/shortener"
)

type createLinkRequestBody struct {
	OriginalURL string `json:"original_url" binding:"required,http_url"`
	ShortName   string `json:"short_name" binding:"omitempty,min=3,max=32"`
}

type linkResponseBody struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

type createLinkResponseBody linkResponseBody
type getLinkResponseBody linkResponseBody
type listLinksResponseBody []linkResponseBody

type updateLinkRequestBody struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name" binding:"required,min=3,max=32"`
}

type updateLinkResponseBody linkResponseBody

func newLinkResponseBody(link shortener.Link, baseURL string) linkResponseBody {
	return linkResponseBody{
		ID:          link.ID,
		OriginalURL: link.OriginalURL,
		ShortName:   link.ShortName,
		ShortURL:    baseURL + "/r/" + link.ShortName,
	}
}

func newListLinksResponseBody(links []shortener.Link, baseURL string) listLinksResponseBody {
	result := make(listLinksResponseBody, len(links))
	for i, link := range links {
		result[i] = newLinkResponseBody(link, baseURL)
	}

	return result
}
