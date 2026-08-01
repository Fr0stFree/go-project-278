package metrics

// LinkVisit represents a record of a link visit, including details.
type LinkVisit struct {
	ID        int    `json:"id"`
	LinkID    int    `json:"link_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Status    int    `json:"status"`
	Referrer  string `json:"referrer"`
}
