package shortener

// Link represents a shortened link with its original URL, short name, and the generated short URL.
type Link struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

// LinkRange represents a range of links to be retrieved, specified by the starting and ending indices.
type LinkRange struct {
	From int
	To   int
}

func (l *LinkRange) IsValid() bool {
	return l.From >= 0 && l.To >= 0 && l.From <= l.To
}
