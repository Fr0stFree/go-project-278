// Package link handles shortened link HTTP requests.
package link

import (
	"fmt"
	"net/http"
	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

type shortenerService interface {
	GetRedirectLink(shortName string) (shortener.Link, error)
	SaveLinkVisit(linkID uint, ip, userAgent, referrer string, status uint) (shortener.LinkVisit, error)
	CreateLink(originalURL, shortName string) (shortener.Link, error)
	GetLink(id uint) (shortener.Link, error)
	ListLinksWithCount(optsBuilder *shortener.LinkListOptionsBuilder) ([]shortener.Link, int, error)
	UpdateLink(id uint, originalURL, shortName string) (shortener.Link, error)
	DeleteLink(id uint) error
}

type handler struct {
	shortener shortenerService
}

func (h *handler) redirect(ctx *gin.Context) {
	shortName := ctx.Param("short_name")

	link, err := h.shortener.GetRedirectLink(shortName)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	referrer := ctx.GetHeader("Referer")
	status := http.StatusFound

	_, err = h.shortener.SaveLinkVisit(link.ID, ip, userAgent, referrer, uint(status))
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	ctx.Redirect(status, link.OriginalURL)
}

func (h *handler) create(ctx *gin.Context) {
	var body createLinkRequestBody

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	link, err := h.shortener.CreateLink(body.OriginalURL, body.ShortName)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, createLinkResponseBody(link))
}

func (h *handler) get(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	link, err := h.shortener.GetLink(linkID)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, getLinkResponseBody(link))
}

func (h *handler) list(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	links, count, err := h.shortener.ListLinksWithCount(optsBuilder)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	from, to := optsBuilder.Range()

	ctx.Header("Content-Range", fmt.Sprintf("links %d-%d/%d", from, to, count))
	ctx.JSON(http.StatusOK, listLinksResponseBody(links))
}

func (h *handler) update(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	var body updateLinkRequestBody

	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	link, err := h.shortener.UpdateLink(linkID, body.OriginalURL, body.ShortName)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, updateLinkResponseBody(link))
}

func (h *handler) delete(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	err = h.shortener.DeleteLink(linkID)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
