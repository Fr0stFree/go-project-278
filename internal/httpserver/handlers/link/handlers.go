// Package link provides HTTP handlers for managing shortened links.
package link

import (
	"fmt"
	"net/http"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

type handler struct {
	Service *shortener.Service
}

func (h *handler) redirect(ctx *gin.Context) {
	shortName := ctx.Param("short_name")

	link, err := h.Service.GetRedirectLink(shortName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.Redirect(http.StatusFound, link.OriginalURL)
}

func (h *handler) create(ctx *gin.Context) {
	var body createLinkRequestBody

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	link, err := h.Service.CreateLink(body.OriginalURL, body.ShortName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, createLinkResponseBody(link))
}

func (h *handler) get(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	link, err := h.Service.GetLink(linkID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, getLinkResponseBody(link))
}

func (h *handler) list(ctx *gin.Context) {
	filterOpts, err := parseFilterOpts(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	links, err := h.Service.ListLinks(filterOpts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	amount, err := h.Service.CountLinks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	from, to := filterOpts.Range()
	ctx.Header("Content-Range", fmt.Sprintf("links %d-%d/%d", from, to, amount))
	ctx.JSON(http.StatusOK, listLinksResponseBody(links))
}

func (h *handler) update(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var body updateLinkRequestBody

	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	link, err := h.Service.UpdateLink(linkID, body.OriginalURL, body.ShortName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, updateLinkResponseBody(link))
}

func (h *handler) delete(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = h.Service.DeleteLink(linkID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.Status(http.StatusNoContent)
}
