// Package metrics provides services for managing link visits and interacting with the link visit repository.
package metrics

import (
	"fmt"
	"net/http"
	"shortener/internal/services/metrics"

	"github.com/gin-gonic/gin"
)

type handler struct {
	metrics *metrics.Service
}

func (h *handler) listLinkVisits(ctx *gin.Context) {
	opts, err := parseFilterOpts(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	visits, err := h.metrics.ListLinkVisits(opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	count, err := h.metrics.CountLinkVisits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	from, to := opts.Range()

	ctx.Header("Content-Range", fmt.Sprintf("link_visits %d-%d/%d", from, to, count))
	ctx.JSON(http.StatusOK, listLinksVisitsResponseBody(visits))
}
