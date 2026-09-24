// Package link handles shortened link HTTP requests.
package link

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/httpserver/httptools/httperror"
	"shortener/internal/httpserver/httptools/httpparam"
	"shortener/internal/services/shortener"
)

// Service defines the link management operations used by the HTTP layer.
type Service interface {
	GetRedirectLink(ctx context.Context, shortName string) (shortener.Link, error)
	SaveLinkVisit(ctx context.Context, linkID uint, ip, userAgent, referrer string, status uint) (shortener.LinkVisit, error)
	CreateLink(ctx context.Context, originalURL, shortName string) (shortener.Link, error)
	GetLink(ctx context.Context, id uint) (shortener.Link, error)
	ListLinksWithCount(ctx context.Context, optsBuilder *shortener.LinkListOptionsBuilder) ([]shortener.Link, int, error)
	UpdateLink(ctx context.Context, id uint, originalURL, shortName string) (shortener.Link, error)
	DeleteLink(ctx context.Context, id uint) error
}

type handler struct {
	service Service
	baseURL string
}

func (h *handler) redirect(ctx *gin.Context) {
	shortName, err := httpparam.ReadStringPath(ctx, "short_name")
	if err != nil {
		httperror.WriteResponse(ctx, shortener.NewValidationError(err.Error(), "short_name"))

		return
	}

	link, err := h.service.GetRedirectLink(ctx.Request.Context(), shortName)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	var (
		ip        = httpparam.ReadClientIP(ctx)
		userAgent = httpparam.ReadUserAgentHeader(ctx)
		referrer  = httpparam.ReadReferrerHeader(ctx)
		status    = http.StatusFound
	)

	_, err = h.service.SaveLinkVisit(ctx.Request.Context(), link.ID, ip, userAgent, referrer, uint(status))
	if err != nil {
		slog.Error(
			"Failed to save link visit; redirect will continue",
			slog.String("request_id", httpparam.ReadRequestIDContext(ctx)),
			slog.Any("error", err),
			slog.Uint64("link_id", uint64(link.ID)),
			slog.String("short_name", shortName),
			slog.Int("status", status),
		)
	}

	ctx.Redirect(status, link.OriginalURL)
}

func (h *handler) create(ctx *gin.Context) {
	var body createLinkRequestBody

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	link, err := h.service.CreateLink(ctx.Request.Context(), body.OriginalURL, body.ShortName)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, createLinkResponseBody(newLinkResponseBody(link, h.baseURL)))
}

func (h *handler) get(ctx *gin.Context) {
	linkID, err := httpparam.ReadNonNegativeIntPath(ctx, "id")
	if err != nil {
		httperror.WriteResponse(ctx, shortener.NewValidationError(err.Error(), "link_id"))

		return
	}

	link, err := h.service.GetLink(ctx.Request.Context(), linkID)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, getLinkResponseBody(newLinkResponseBody(link, h.baseURL)))
}

func (h *handler) list(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	links, count, err := h.service.ListLinksWithCount(ctx.Request.Context(), optsBuilder)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	from, _ := optsBuilder.Range()
	httpparam.WriteContentRangeHeader(ctx, "links", from, len(links), count)
	ctx.JSON(http.StatusOK, newListLinksResponseBody(links, h.baseURL))
}

func (h *handler) update(ctx *gin.Context) {
	linkID, err := httpparam.ReadNonNegativeIntPath(ctx, "id")
	if err != nil {
		httperror.WriteResponse(ctx, shortener.NewValidationError(err.Error(), "link_id"))

		return
	}

	var body updateLinkRequestBody

	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	link, err := h.service.UpdateLink(ctx.Request.Context(), linkID, body.OriginalURL, body.ShortName)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	ctx.JSON(http.StatusOK, updateLinkResponseBody(newLinkResponseBody(link, h.baseURL)))
}

func (h *handler) delete(ctx *gin.Context) {
	linkID, err := httpparam.ReadNonNegativeIntPath(ctx, "id")
	if err != nil {
		httperror.WriteResponse(ctx, shortener.NewValidationError(err.Error(), "link_id"))

		return
	}

	err = h.service.DeleteLink(ctx.Request.Context(), linkID)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkListOptionsBuilder, error) {
	builder := shortener.NewLinkListOptionsBuilder()

	rangeQuery, err := httpparam.ReadRangeQuery(ctx, "range")
	if err != nil {
		return nil, shortener.NewValidationError(err.Error(), "range")
	}

	if rangeQuery != nil {
		builder.WithRange(rangeQuery.From, rangeQuery.Count)
	}

	sortQuery, err := httpparam.ReadSortQuery(ctx, "sort")
	if err != nil {
		return nil, shortener.NewValidationError(err.Error(), "sort")
	}

	if sortQuery != nil {
		sortField := sortQuery.Field
		if sortField == "short_url" {
			sortField = string(shortener.LinkSortByShortName)
		}

		builder.WithSort(sortField, sortQuery.Direction)
	}

	return builder, nil
}
