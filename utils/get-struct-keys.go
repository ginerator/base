package utils

import (
	"reflect"
	"strings"
)

func GetStructKeys(s interface{}) []string {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return nil
	}

	var keys []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := strings.ToLower(string(field.Name[0])) + field.Name[1:]

		if field.Type.Kind() == reflect.Struct {
			nestedKeys := GetStructKeys(reflect.New(field.Type).Elem().Interface())
			for _, nestedKey := range nestedKeys {
				keys = append(keys, nestedKey)
			}
		} else {
			keys = append(keys, fieldName)
		}
	}
	return keys
}
