//go:build unit

package user_test

import (
	"encoding/json"
	"fmt"
	"testing"

	user "github.com/ginerator/base/model/users"
	"github.com/stretchr/testify/assert"
)

type UserTypeStruct struct {
	Type user.UserType
}

func TestUserTypeEnumUnmarshalJSON(t *testing.T) {
	type testData struct {
		UserTypeStructJSON string
		expected           user.UserType
	}

	testCases := make([]testData, 0)
	for _, userType := range user.UserTypes {
		testCases = append(testCases, testData{
			UserTypeStructJSON: fmt.Sprintf(`{"type": "%s"}`, userType),
			expected:           userType,
		})
	}

	for _, testCase := range testCases {
		var userTypeStruct UserTypeStruct
		err := json.Unmarshal([]byte(testCase.UserTypeStructJSON), &userTypeStruct)
		assert.NoError(t, err)
		actual := userTypeStruct.Type
		assert.Equal(t, actual, testCase.expected)
	}
}
