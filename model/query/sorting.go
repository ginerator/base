package query

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc               = "DESC"
)

type Sorting interface {
	GetSort() SortDirection
	GetSortBy() string
	GetSortableAttributes() map[string]string
}
