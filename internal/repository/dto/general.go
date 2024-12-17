package dto

type FilterOptions struct {
	Filters map[string]any
	SortBy  string
	Order   string
	Page    int
	Limit   int
}
