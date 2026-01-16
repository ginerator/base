package postgres

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ginerator/base/errors"
	modelquery "github.com/ginerator/base/model/query"
	"github.com/ginerator/base/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func Create[R interface{}, M interface{}](ctx *gin.Context, client *BunPostgresDatabaseClient, createItemRequest R) (M, error) {
	db := client.getDB(ctx)

	entity := new(M)
	_, err := db.NewInsert().Model(&createItemRequest).Returning("*").Exec(ctx, entity)
	if err != nil {
		log.Error().
			Err(err).
			Str("model", fmt.Sprintf("%T", *entity)).
			Msg("[BASE REPOSITORY] - Create - Inserting new entity")
		return *entity, errors.NewUnkownDatabaseError(err)
	}
	return *entity, nil
}

func GetOne[M interface{}](ctx *gin.Context, client *BunPostgresDatabaseClient, id uuid.UUID, userId *string) (M, error) {
	db := client.getDB(ctx)

	entity := new(M)
	query := db.NewSelect().Model(entity).Where("id = ?", id).Where("deleted_at IS NULL")
	if userId != nil {
		query.Where("userId = ?", userId)
	}

	err := query.Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().
				Err(err).
				Str("id", id.String()).
				Str("model", fmt.Sprintf("%T", *entity)).
				Msg("[BASE REPOSITORY] - GetOne - Not found")
			return *entity, errors.NewNotFoundError("NOT_FOUND", fmt.Errorf("Entity with id %s could not be found.", id))
		}
		log.Error().
			Err(err).
			Str("id", id.String()).
			Str("model", fmt.Sprintf("%T", *entity)).
			Msg("[BASE REPOSITORY] - GetOne - Unhandled error")
		return *entity, errors.NewInternalServerError("UNKNOWN_ERROR", err)
	}

	return *entity, nil
}

func GetMany[Q interface{}, M interface{}](ctx *gin.Context, client *BunPostgresDatabaseClient, query Q, userId *string) ([]M, modelquery.ResponseMeta, error) {
	db := client.getDB(ctx)
	entities := make([]M, 0)
	entity := new(M) // Just to show it in a log
	responseMeta := modelquery.ResponseMeta{}

	dbQuery := db.NewSelect().Model(&entities)
	if userId != nil {
		dbQuery.Where("userId = ?", userId)
	}

	offset, limit := utils.BuildQuery(ctx, dbQuery)

	count, err := dbQuery.ScanAndCount(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().
				Err(err).
				Msg("[BASE REPOSITORY] - GetMany - Not found")
			return entities, responseMeta, nil
		}
		log.Error().
			Err(err).
			Str("model", fmt.Sprintf("%T", *entity)).
			Msg("[BASE REPOSITORY] - GetMany - Unhandled error")
		return entities, responseMeta, errors.NewInternalServerError("UNKNOWN_ERROR", err)
	}

	return entities, utils.BuildResponseMeta(offset, limit, count), nil
}

func UpdateOne[M interface{}, R interface{}](ctx *gin.Context, client *BunPostgresDatabaseClient, id uuid.UUID, request R, userId *string) (M, error) {
	entity := new(M)

	query := client.getDB(ctx).NewUpdate().OmitZero().Model(&request).Where("id = ?", id).Where("deleted_at IS NULL").Returning("*")
	if userId != nil {
		log.Debug().
			Str("id", id.String()).
			Str("model", fmt.Sprintf("%T", *entity)).
			Msg("[BASE REPOSITORY] - UpdateOne - Fetching with userId")
		query.Where("userId = ?", userId)
	}

	_, err := query.Exec(ctx, entity)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", id.String()).
			Str("model", fmt.Sprintf("%T", *entity)).
			Msg("[BASE REPOSITORY] - UpdateOne - Error updating")
		return *entity, err
	}
	return *entity, nil
}
