package users

import "strconv"

type UserAccess int

const (
	UANone UserAccess = iota
	UARevoked
	UAMember
	UAAdmin
)

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Access    UserAccess `json:"access"`
}

func (u User) IsAdmin() bool {
	return u.Access == UAAdmin
}

func (u User) IsMember() bool {
	return u.Access == UAMember
}

func (u User) IsRevoked() bool {
	return u.Access == UARevoked
}

func (u User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	} else {
		if u.LastName != "" {
			return u.FirstName + " " + u.LastName
		}
		return u.FirstName
	}
}

func (u User) Recipient() string {
	return strconv.Itoa(u.ID)
}
