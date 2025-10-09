package models

import "github.com/google/uuid"

// Member represents a lightweight user in project context
type Member struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// ToMember converts a User to a Member
func (u *User) ToMember() Member {
	return Member{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

// ToMembers converts a slice of Users to Members
func ToMembers(users []User) []Member {
	members := make([]Member, len(users))
	for i, user := range users {
		members[i] = user.ToMember()
	}
	return members
}
