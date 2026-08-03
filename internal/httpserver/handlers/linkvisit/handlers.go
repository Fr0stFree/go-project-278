package linkvisit

import (
	"fmt"
	"net/http"

	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

type handler struct {
	shortener *shortener.Service
}

func (h *handler) listLinkVisits(ctx *gin.Context) {
	optsBuilder, err := parseFilterOpts(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	visits, count, err := h.shortener.ListLinkVisitsWithCount(optsBuilder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	from, to := optsBuilder.Range()

	ctx.Header("Content-Range", fmt.Sprintf("link_visits %d-%d/%d", from, to, count))
	ctx.JSON(http.StatusOK, listLinksVisitsResponseBody(visits))
}
