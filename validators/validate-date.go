package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func IsValidDate(fl validator.FieldLevel) bool {
	exampleDate := "2010-10-20"
	layoutDate := "2006-01-02"
	_, err := time.Parse(layoutDate, exampleDate)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing string date")
		return false
	}
	return true
}
