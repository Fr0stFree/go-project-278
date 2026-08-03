package linkvisit

import (
	"fmt"
	"net/http"

	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

type handler struct {
	shortener *shortener.Service
}

func (h *handler) listLinkVisits(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	visits, count, err := h.shortener.ListLinkVisitsWithCount(optsBuilder)
	if err != nil {
		httpserver.WriteErrorResponse(ctx, err)

		return
	}

	from, to := optsBuilder.Range()

	ctx.Header("Content-Range", fmt.Sprintf("link_visits %d-%d/%d", from, to, count))
	ctx.JSON(http.StatusOK, listLinksVisitsResponseBody(visits))
}
