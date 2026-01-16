package utils

import (
	modelquery "github.com/PlanToPack/api-utils/model/query"
)

func BuildResponseMeta(offset int, limit int, count int) modelquery.ResponseMeta {
	return modelquery.ResponseMeta{
		ItemsTotal: count,
		PagesTotal: (int(count/limit) + 1),
		Pagination: modelquery.Pagination{
			PageSize: limit,
			Page:     offset,
		},
	}
}
