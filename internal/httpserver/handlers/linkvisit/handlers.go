package linkvisit

import (
	"context"
	"fmt"
	"net/http"
	"shortener/internal/httpserver/httptools"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

type shortenerService interface {
	ListLinkVisitsWithCount(ctx context.Context, optsBuilder *shortener.LinkVisitListOptionsBuilder) ([]shortener.LinkVisit, int, error)
}
type handler struct {
	shortener shortenerService
}

func (h *handler) list(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	visits, count, err := h.shortener.ListLinkVisitsWithCount(ctx.Request.Context(), optsBuilder)
	if err != nil {
		httptools.WriteErrorResponse(ctx, err)

		return
	}

	from, to := optsBuilder.Range()

	ctx.Header("Content-Range", fmt.Sprintf("link_visits %d-%d/%d", from, to, count))
	ctx.JSON(http.StatusOK, listLinksVisitsResponseBody(visits))
}
