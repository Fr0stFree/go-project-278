package link

import (
	"net/http"
	"shortener/internal/services/shortener"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *shortener.Service
}

func RegisterRoutes(service *shortener.Service, router gin.IRouter) {
	h := &Handler{Service: service}
	router.POST("/api/links", h.Create)
	router.GET("/api/links/:id", h.Get)
	router.GET("/api/links", h.List)
	router.PUT("/api/links/:id", h.Update)
	router.DELETE("/api/links/:id", h.Delete)
}

// Create generates a short URL for the given original URL and returns it in the response.
func (l *Handler) Create(ctx *gin.Context) {
	var body createLinkRequestBody

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	link, err := l.Service.CreateLink(body.OriginalURL, body.ShortName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, createLinkResponseBody(link))
}

// Get retrieves the original URL corresponding to the given short URL.
func (l *Handler) Get(ctx *gin.Context) {
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
func (l *Handler) List(ctx *gin.Context) {
	var linkRange *shortener.LinkRange
	linkRange, err := parseLinksRange(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	links, err := l.Service.ListLinks(linkRange)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, listLinksResponseBody(links))
}

func (l *Handler) Update(ctx *gin.Context) {
	linkID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})

		return
	}

	var body updateLinkRequestBody

	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	link, err := l.Service.UpdateLink(linkID, body.OriginalURL, body.ShortName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, updateLinkResponseBody(link))
}

func (l *Handler) Delete(ctx *gin.Context) {
	linkID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})

		return
	}

	err = l.Service.DeleteLink(linkID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.Status(http.StatusNoContent)
}
