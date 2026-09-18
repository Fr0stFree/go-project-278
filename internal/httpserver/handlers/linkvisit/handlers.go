package linkvisit

import (
	"context"
	"net/http"
	"shortener/internal/httpserver/httptools/httperror"
	"shortener/internal/httpserver/httptools/httpparam"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

// Service defines the link management operations used by the HTTP layer.
type Service interface {
	ListLinkVisitsWithCount(ctx context.Context, optsBuilder *shortener.LinkVisitListOptionsBuilder) ([]shortener.LinkVisit, int, error)
}

type handler struct {
	service Service
}

func (h *handler) list(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	visits, count, err := h.service.ListLinkVisitsWithCount(ctx.Request.Context(), optsBuilder)
	if err != nil {
		httperror.WriteResponse(ctx, err)

		return
	}

	from, _ := optsBuilder.Range()
	httpparam.WriteContentRangeHeader(ctx, "link_visits", from, len(visits), count)
	ctx.JSON(http.StatusOK, listLinksVisitsResponseBody(visits))
}

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkVisitListOptionsBuilder, error) {
	builder := shortener.NewLinkVisitListOptionsBuilder()

	rangeQuery, err := httpparam.ReadRangeQuery(ctx, "range")
	if err != nil {
		return nil, shortener.NewValidationError(err.Error(), "range")
	}

	if rangeQuery != nil {
		builder.WithRange(rangeQuery.From, rangeQuery.Count)
	}

	return builder, nil
}
