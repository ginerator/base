package utils

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func UnmarshalJSONGeneric[T comparable](byteArray []byte, validValues []T) (T, error) {
	var zeroValue T
	str := strings.Trim(string(byteArray), `"`)

	upperCasedValue := strings.ToUpper(str)

	for _, valid := range validValues {
		if strings.EqualFold(fmt.Sprintf("%v", valid), upperCasedValue) {
			return valid, nil
		}
	}

	stringValuesSupported := lo.Reduce(validValues, func(agg string, item T, _ int) string {
		if agg == "" {
			return fmt.Sprintf("%v", item)
		}
		return fmt.Sprintf("%s, %v", agg, item)
	}, "")

	return zeroValue, fmt.Errorf("Value '%v' isn't a valid value. The valid values are: %s", str, stringValuesSupported)
}
