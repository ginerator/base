package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PlanToPack/api-utils/errors"
	"github.com/PlanToPack/api-utils/model/query"
	"github.com/PlanToPack/api-utils/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func formatValidationErrors(errs error) *errors.CustomError {
	if _, ok := errs.(*validator.InvalidValidationError); ok {
		return errors.NewInvalidPayloadError("INVALID_PAYLOAD", errs)
	}

	err := errs.(validator.ValidationErrors)[0]

	return errors.NewInvalidPayloadError("INVALID_PAYLOAD", fmt.Errorf("Value '%s' for attribute '%s' is not of type: %s", err.Value(), err.Field(), err.Tag()))
}

func Create[R interface{}, M interface{}](ctx *gin.Context, validator *validator.Validate, serviceFunction func(*gin.Context, R) (M, error)) {
	var request R

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - Create - Error decoding request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors := validator.Struct(request)
	if validationErrors != nil {
		log.Error().Err(validationErrors).Msg("[BASE CONTROLLER] - Create - Error validating struct")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": formatValidationErrors(validationErrors).Error()})
		return
	}

	entity, err := serviceFunction(ctx, request)
	if err != nil {
		log.Error().Err(validationErrors).Msg("[BASE CONTROLLER] - Create - Error in service function")
		ctx.JSON(400, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": entity})
}

func CreateWithExternalId[R interface{}, M interface{}](ctx *gin.Context, validator *validator.Validate, serviceFunction func(*gin.Context, uuid.UUID, R) (M, error)) {
	externalId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var request R

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - CreateWithExternalId - Error decoding request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors := validator.Struct(request)
	if validationErrors != nil {
		log.Error().Err(validationErrors).Msg("[BASE CONTROLLER] - CreateWithExternalId - Error validating struct")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": formatValidationErrors(validationErrors).Error()})
		return
	}

	entity, err := serviceFunction(ctx, externalId, request)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - CreateWithExternalId - Error in service function")
		ctx.JSON(400, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": entity})
}

func validateQuery[Q interface{}](ctx *gin.Context, allowedQueryParams []string) (Q, error) {
	// Filter out unkown fields
	unknowFields, _ := lo.Difference(
		lo.Keys[string, []string](ctx.Request.URL.Query()),
		allowedQueryParams,
	)

	var query Q
	if len(unknowFields) > 0 {
		err := fmt.Errorf("The following field(s) are not allowed: %s. Allowed fields are: %s", strings.Join(unknowFields, ", "), strings.Join(allowedQueryParams, ", "))
		log.Error().Err(err).Msg("[BASE CONTROLLER] - validateQuery - Unknown fields")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return query, err
	}

	// Parse query
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - validateQuery - Does not bind query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return query, err
	}

	return query, nil
}

func GetOne[M interface{}](ctx *gin.Context, serviceFunction func(*gin.Context, uuid.UUID) (M, error)) {
	uuid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetOne - Retrieving id")
		ctx.Error(err)
		return
	}

	entity, err := serviceFunction(ctx, uuid)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetOne - Calling service function")
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": entity})
}

func GetOneHidrated[Q interface{}, M interface{}, A interface{}](ctx *gin.Context, allowedQueryParams []string, queryParser func(Q) (A, error), serviceFunction func(*gin.Context, uuid.UUID, A) (M, error)) {
	uuid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetOneHydrated - Retrieving id")
		ctx.Error(err)
		return
	}

	query, err := validateQuery[Q](ctx, allowedQueryParams)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetOneHydrated - Validating query")
		ctx.Error(err)
		return
	}

	var additionalQueryData A
	if queryParser != nil {
		additionalQueryData, err = queryParser(query)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	entity, err := serviceFunction(ctx, uuid, additionalQueryData)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetOneHydrated - Calling service function")
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": entity})
}

func GetMany[Q interface{}, M interface{}](ctx *gin.Context, validator *validator.Validate, serviceFunction func(*gin.Context, Q) ([]M, query.ResponseMeta, error)) {
	var query Q
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetMany - Does not bind query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValidQuery, unknownFields := validators.IsValidQuery(ctx, query)
	if !isValidQuery {
		log.Error().Str("unknownFields", strings.Join(unknownFields, ", ")).Msg("[BASE CONTROLLER] - GetMany - Invalid query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query params"})
		return

	}

	validationErrors := validator.Struct(query)
	if validationErrors != nil {
		log.Error().Err(validationErrors).Msg("[BASE CONTROLLER] - GetMany - Validating struct")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": formatValidationErrors(validationErrors).Error()})
		return
	}

	entitys, responseMeta, err := serviceFunction(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetMany - Calling service function")
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"meta": responseMeta,
		"data": entitys,
	})
}

func GetManyWithExternalId[Q interface{}, M interface{}](ctx *gin.Context, validator *validator.Validate, serviceFunction func(*gin.Context, uuid.UUID, Q) ([]M, query.ResponseMeta, error)) {
	externalId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetManyWithExternalId - Retrieving id")
		ctx.Error(err)
		return
	}

	var query Q
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetManyWithExternalId - Does not bind query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValidQuery, unknownFields := validators.IsValidQuery(ctx, query)
	if !isValidQuery {
		log.Error().Str("unknownFields", strings.Join(unknownFields, ", ")).Msg("[BASE CONTROLLER] - GetManyWithExternalId - Invalid query params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query params"})
		return

	}

	validationErrors := validator.Struct(query)
	if validationErrors != nil {
		log.Error().Err(validationErrors).Msg("[BASE CONTROLLER] - GetManyWithExternalId - Validating struct")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": formatValidationErrors(validationErrors).Error()})
		return
	}

	entitys, responseMeta, err := serviceFunction(ctx, externalId, query)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - GetManyWithExternalId - Calling service function")
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"meta": responseMeta,
		"data": entitys,
	})
}

func UpdateOne[R interface{}, M interface{}](ctx *gin.Context, validator *validator.Validate, serviceFunction func(*gin.Context, uuid.UUID, R) (M, error)) {
	uuid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - UpdateOne - Retrieving id")
		ctx.Error(err)
		return
	}
	var request R

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - UpdateOne - Decoding struct")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors := validator.Struct(request)
	if validationErrors != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - UpdateOne - Validating struct")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": formatValidationErrors(validationErrors).Error()})
		return
	}

	entityUpdated, err := serviceFunction(ctx, uuid, request)
	if err != nil {
		log.Error().Err(err).Msg("[BASE CONTROLLER] - UpdateOne - Calling service function")
		ctx.JSON(400, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": entityUpdated})
}
