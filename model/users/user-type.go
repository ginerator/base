package user

import "github.com/PlanToPack/api-utils/utils"

type UserType string

const (
	UserTypeSystem UserType = "SYSTEM"
	UserTypePerson UserType = "PERSON"
)

var UserTypes = []UserType{
	UserTypeSystem,
	UserTypePerson,
}

func (userType *UserType) UnmarshalJSON(byteArray []byte) error {
	value, err := utils.UnmarshalJSONGeneric(byteArray, UserTypes)
	if err != nil {
		return err
	}
	*userType = value
	return nil
}
