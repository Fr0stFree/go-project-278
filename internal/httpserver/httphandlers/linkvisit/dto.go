// Package linkvisit handles redirect visit HTTP requests.
package linkvisit

import (
	"shortener/internal/services/shortener"
	"time"
)

type linkVisitResponseBody struct {
	ID        uint   `json:"id"`
	LinkID    uint   `json:"link_id"`
	CreatedAt string `json:"created_at"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Status    uint   `json:"status"`
}

type listLinksVisitsResponseBody []linkVisitResponseBody

func newListLinkVisitsResponseBody(visits []shortener.LinkVisit) listLinksVisitsResponseBody {
	result := make(listLinksVisitsResponseBody, len(visits))
	for i, visit := range visits {
		result[i] = linkVisitResponseBody{
			ID:        visit.ID,
			LinkID:    visit.LinkID,
			CreatedAt: visit.CreatedAt.Format(time.RFC3339),
			IP:        visit.IP,
			UserAgent: visit.UserAgent,
			Status:    visit.Status,
		}
	}

	return result
}
