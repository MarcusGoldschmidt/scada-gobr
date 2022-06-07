package models

import "github.com/google/uuid"

type UserView struct {
	Id     uuid.UUID `json:"id" gorm:"type:uuid"`
	ViewId uuid.UUID `json:"viewId" gorm:"type:uuid"`
	UserId uuid.UUID `json:"userId" gorm:"type:uuid"`
}
