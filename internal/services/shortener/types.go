package shortener

// Link represents a shortened link with its original URL, short name, and the generated short URL.
type Link struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}
