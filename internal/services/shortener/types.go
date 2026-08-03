package shortener

// Link represents a shortened link with its original URL, short name, and the generated short URL.
type Link struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

// LinkVisit represents a record of a link visit, including details.
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
