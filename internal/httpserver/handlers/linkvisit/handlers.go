package linkvisit

import (
	"context"
	"net/http"
	"shortener/internal/httpserver/httptools"
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
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	visits, count, err := h.service.ListLinkVisitsWithCount(ctx.Request.Context(), optsBuilder)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	from, _ := optsBuilder.Range()
	ctx.Header("Content-Range", httptools.FormatContentRange("link_visits", from, len(visits), count))
	ctx.JSON(http.StatusOK, listLinksVisitsResponseBody(visits))
}

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkVisitListOptionsBuilder, error) {
	builder := shortener.NewLinkVisitListOptionsBuilder()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		from, to, err := httptools.ParseQueryRange(rangeRaw)
		if err != nil {
			return nil, shortener.NewValidationError(err.Error(), "range")
		}

		builder.WithRange(from, to)
	}

	return builder, nil
}
