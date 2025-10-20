package models

import "github.com/google/uuid"

type Member struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (u *User) ToMember() Member {
	return Member{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToMembers(users []User) []Member {
	members := make([]Member, len(users))
	for i, user := range users {
		members[i] = user.ToMember()
	}
	return members
}
