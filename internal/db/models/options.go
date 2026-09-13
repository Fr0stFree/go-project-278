package models

// ListOptions stores pagination and sorting values shared by storage queries.
type ListOptions struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}
