package link

import (
	"shortener/internal/services/shortener"
)

type createLinkRequestBody struct {
	OriginalURL string `json:"original_url" binding:"required,http_url"`
	ShortName   string `json:"short_name" binding:"omitempty,excludes=/"`
}

type createLinkResponseBody shortener.Link

type getLinkResponseBody shortener.Link

type listLinksResponseBody []shortener.Link

type updateLinkRequestBody struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name" binding:"required,excludes=/"`
}

type updateLinkResponseBody shortener.Link
