// Package link handles shortened link HTTP requests.
package link

import (
	"context"
	"fmt"
	"net/http"
	"shortener/internal/httpserver/httptools"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

type shortenerService interface {
	GetRedirectLink(ctx context.Context, shortName string) (shortener.Link, error)
	SaveLinkVisit(ctx context.Context, linkID uint, ip, userAgent, referrer string, status uint) (shortener.LinkVisit, error)
	CreateLink(ctx context.Context, originalURL, shortName string) (shortener.Link, error)
	GetLink(ctx context.Context, id uint) (shortener.Link, error)
	ListLinksWithCount(ctx context.Context, optsBuilder *shortener.LinkListOptionsBuilder) ([]shortener.Link, int, error)
	UpdateLink(ctx context.Context, id uint, originalURL, shortName string) (shortener.Link, error)
	DeleteLink(ctx context.Context, id uint) error
}

type handler struct {
	shortener shortenerService
}

func (h *handler) redirect(ctx *gin.Context) {
	shortName := ctx.Param("short_name")

	link, err := h.shortener.GetRedirectLink(ctx.Request.Context(), shortName)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	referrer := ctx.GetHeader("Referer")
	status := http.StatusFound

	_, err = h.shortener.SaveLinkVisit(ctx.Request.Context(), link.ID, ip, userAgent, referrer, uint(status))
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	ctx.Redirect(status, link.OriginalURL)
}

func (h *handler) create(ctx *gin.Context) {
	var body createLinkRequestBody

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	link, err := h.shortener.CreateLink(ctx.Request.Context(), body.OriginalURL, body.ShortName)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, createLinkResponseBody(link))
}

func (h *handler) get(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	link, err := h.shortener.GetLink(ctx.Request.Context(), linkID)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, getLinkResponseBody(link))
}

func (h *handler) list(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	links, count, err := h.shortener.ListLinksWithCount(ctx.Request.Context(), optsBuilder)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	from, to := optsBuilder.Range()

	ctx.Header("Content-Range", fmt.Sprintf("links %d-%d/%d", from, to, count))
	ctx.JSON(http.StatusOK, listLinksResponseBody(links))
}

func (h *handler) update(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	var body updateLinkRequestBody

	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	link, err := h.shortener.UpdateLink(ctx.Request.Context(), linkID, body.OriginalURL, body.ShortName)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, updateLinkResponseBody(link))
}

func (h *handler) delete(ctx *gin.Context) {
	linkID, err := parseLinkID(ctx)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	err = h.shortener.DeleteLink(ctx.Request.Context(), linkID)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

func parseLinkID(ctx *gin.Context) (uint, error) {
	linkID, err := httptools.ParsePositiveIntParam(ctx.Param("id"))
	if err != nil {
		return 0, shortener.NewValidationError(err.Error(), "link_id")
	}

	return uint(linkID), nil
}

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkListOptionsBuilder, error) {
	builder := shortener.NewLinkListOptionsBuilder()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		from, to, err := httptools.ParseQueryRange(rangeRaw)
		if err != nil {
			return nil, shortener.NewValidationError(err.Error(), "range")
		}

		builder.WithRange(from, to)
	}

	sortRaw := ctx.Query("sort")
	if sortRaw != "" {
		sortBy, sortOrder, err := httptools.ParseQuerySort(sortRaw)
		if err != nil {
			return nil, shortener.NewValidationError(err.Error(), "sort")
		}

		builder.WithSort(sortBy, sortOrder)
	}

	return builder, nil
}
