// Package linkvisit handles redirect visit HTTP requests.
package linkvisit

import (
	"shortener/internal/services/shortener"
)

type listLinksVisitsResponseBody []shortener.LinkVisit
