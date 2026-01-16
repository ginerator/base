package user

type User struct {
	Email  *string  `json:"email"`
	Id     *string  `json:"id"`
	Source *string  `json:"source"`
	Type   UserType `json:"type"`
}
