package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

var defaultQueryControlParams = map[string]string{
	"sortBy": "created_at",
	"sort":   "DESC",
	"limit":  "10",
	"offset": "0",
}

var queryControlParams = func() []string {
	params := make([]string, 0, len(defaultQueryControlParams))
	for param := range defaultQueryControlParams {
		params = append(params, param)
	}
	return params
}()

func filterOutDeletedEntities(dbQuery *bun.SelectQuery) {
	dbQuery.Where("deleted_at IS NULL")
}

func urlToDbQuery(gCtx *gin.Context, dbQuery *bun.SelectQuery) {
	for param, values := range gCtx.Request.URL.Query() {
		if !lo.Contains(queryControlParams, param) {
			if len(values) > 1 {
				dbQuery.Where(fmt.Sprintf("%s IN (?)", param), bun.In(values))
			} else if len(values) == 1 {
				dbQuery.Where(fmt.Sprintf("%s = ?", param), strcase.ToSnake(values[0]))
			}
		}
	}
}

func setQueryControlParams(gCtx *gin.Context, dbQuery *bun.SelectQuery) (int, int) {
	offset, _ := strconv.Atoi(defaultQueryControlParams["offset"])
	limit, _ := strconv.Atoi(defaultQueryControlParams["limit"])

	if sortBy, exists := gCtx.GetQuery("sortBy"); exists {
		if sort, exists := gCtx.GetQuery("sort"); exists {
			dbQuery.Order(fmt.Sprintf("%s %s", strcase.ToSnake(sortBy), sort))
		} else {
			dbQuery.Order(fmt.Sprintf("%s %s", strcase.ToSnake(sortBy), defaultQueryControlParams["sort"]))
		}
	} else {
		dbQuery.Order(fmt.Sprintf("%s %s", strcase.ToSnake(defaultQueryControlParams["sortBy"]), defaultQueryControlParams["sort"]))
	}

	if queryLimit, exists := gCtx.GetQuery("limit"); exists {
		dbQuery.Limit(gCtx.GetInt(queryLimit))
		limit = gCtx.GetInt(queryLimit)
	} else {
		dbQuery.Limit(limit)
	}

	if queryOffset, exists := gCtx.GetQuery("offset"); exists {
		dbQuery.Offset(gCtx.GetInt(queryOffset))
		offset = gCtx.GetInt(queryOffset)
	} else {
		dbQuery.Offset(offset)
	}

	return offset, limit
}

func BuildQuery(gCtx *gin.Context, dbQuery *bun.SelectQuery) (int, int) {
	urlToDbQuery(gCtx, dbQuery)
	filterOutDeletedEntities(dbQuery)
	return setQueryControlParams(gCtx, dbQuery)
}
