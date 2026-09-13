package shortener

import (
	"context"
	"shortener/internal/db/models/link"
	"shortener/internal/db/models/linkvisit"
)

// Link is the service response model for a shortened link.
type Link struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

// LinkVisit is the service response model for a redirect visit.
type LinkVisit struct {
	ID        uint   `json:"id"`
	LinkID    uint   `json:"link_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"-"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Status    uint   `json:"status"`
	Referrer  string `json:"-"`
}

// LinkRepository describes storage operations required for links.
type LinkRepository interface {
	CreateOne(ctx context.Context, insert link.Insert) (link.Record, error)
	GetByID(ctx context.Context, ID uint) (link.Record, error)
	GetMany(ctx context.Context, options link.ListOptions) ([]link.Record, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, ID uint, update link.Update) (link.Record, error)
	DeleteByID(ctx context.Context, ID uint) error
}

// LinkVisitRepository describes storage operations required for visits.
type LinkVisitRepository interface {
	CreateOne(ctx context.Context, insert linkvisit.Insert) (linkvisit.Record, error)
	GetMany(ctx context.Context, options linkvisit.ListOptions) ([]linkvisit.Record, error)
	Count(ctx context.Context) (int, error)
}
