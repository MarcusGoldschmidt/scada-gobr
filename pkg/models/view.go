package models

import (
	"github.com/google/uuid"
	"time"
)

//go:generate stringer -type=ViewType -output=view_type_string.go
type ViewType int

const (
	TimeSeriesViewType ViewType = iota
	GraphicalViewType
	TextViewType
)

type View struct {
	Id             uuid.UUID        `json:"id" gorm:"type:uuid"`
	Name           string           `json:"name"`
	ViewComponents []*ViewComponent `json:"viewComponents" gorm:"foreignKey:ViewId;constraint:OnDelete:CASCADE" `
}

type ViewComponent struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid"`
	Position Position  `json:"position" gorm:"embedded"`
	ViewType ViewType  `json:"viewType"`
	ViewId   uuid.UUID `json:"viewId" gorm:"type:uuid"`
	Data     JSONB     `json:"data" gorm:"type:jsonb"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPosition(x int, y int) *Position {
	return &Position{x, y}
}

type TimeSeriesView struct {
	DataPointsId []uuid.UUID   `json:"dataPointsId"`
	Period       time.Duration `json:"period"`
	Width        int           `json:"width"`
}

type GraphicalView struct {
	Id     uuid.UUID `json:"id" gorm:"type:uuid"`
	Images [][]byte  `json:"images"`
	Min    float64   `json:"min"`
	Max    float64   `json:"max"`
}

type TextView struct {
	Text string `json:"text"`
}
