package shortener

import "time"

// Link is a shortened link used by the application layer.
type Link struct {
	ID          uint
	OriginalURL string
	ShortName   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateLinkParams contains the data needed to persist a new link.
type CreateLinkParams struct {
	OriginalURL string
	ShortName   string
}

// UpdateLinkParams contains the mutable link fields.
type UpdateLinkParams struct {
	OriginalURL string
	ShortName   string
}

// LinkVisit is a redirect visit used by the application layer.
type LinkVisit struct {
	ID        uint
	LinkID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
	IP        string
	UserAgent string
	Status    uint
	Referrer  string
}

// CreateLinkVisitParams contains the data needed to persist a visit.
type CreateLinkVisitParams struct {
	LinkID    uint
	IP        string
	UserAgent string
	Status    uint
	Referrer  string
}
