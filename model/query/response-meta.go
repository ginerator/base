package query

type ResponseMeta struct {
	Pagination
	ItemsTotal int `json:"itemsTotal"`
	PagesTotal int `json:"pagesTotal"`
}
