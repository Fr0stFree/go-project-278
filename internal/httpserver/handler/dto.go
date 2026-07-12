package handler

import (
	"shortener/internal/shortener"
)

type createLinkRequestBody struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name"`
}

type createLinkResponseBody shortener.Link

type getLinkResponseBody shortener.Link

// type listLinksResponseBody []shortener.Link

// type updateLinkRequestBody struct {
// 	OriginalURL string `json:"original_url" binding:"required"`
// 	ShortName   string `json:"short_name" binding:"required"`
// }

// type updateLinkResponseBody shortener.Link
