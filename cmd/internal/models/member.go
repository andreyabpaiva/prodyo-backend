package models

import "github.com/google/uuid"

type Member struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ProjectID uuid.UUID `json:"project_id,omitempty"`
}

func (u *User) ToMember() Member {
	var projectID uuid.UUID
	if u.ProjectID != nil {
		projectID = *u.ProjectID
	}
	return Member{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		ProjectID: projectID,
	}
}

func ToMembers(users []User) []Member {
	members := make([]Member, len(users))
	for i, user := range users {
		members[i] = user.ToMember()
	}
	return members
}
