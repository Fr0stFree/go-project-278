package handler

import (
	"net/http"
	"shortener/internal/shortener"
	"strconv"

	"github.com/gin-gonic/gin"
)

type linkHandler struct {
	Service *shortener.Service
}

func newLinkHandler(service *shortener.Service) *linkHandler {
	return &linkHandler{
		Service: service,
	}
}

// Create generates a short URL for the given original URL and returns it in the response.
func (l *linkHandler) Create(ctx *gin.Context) {
	var body createLinkRequestBody

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	link, err := l.Service.ShortenLink(body.OriginalURL, body.ShortName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, createLinkResponseBody(link))
}

// Get retrieves the original URL corresponding to the given short URL.
func (l *linkHandler) Get(ctx *gin.Context) {
	linkID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
	}

	link, err := l.Service.GetLink(linkID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, getLinkResponseBody(link))
}

// List retrieves a list of all shortened links.
func (l *linkHandler) List(ctx *gin.Context) {
	links, err := l.Service.ListLinks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, listLinksResponseBody(links))
}
