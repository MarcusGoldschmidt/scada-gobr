package postgres

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"time"
)

type ErrorAgg struct {
	Error        string
	CreationTime time.Time
}

type QueueMessages struct {
	Id          uuid.UUID                   `gorm:"type:uuid;primary_key"`
	Headers     shared.JsonB[string]        `gorm:"type:jsonb"`
	QueueName   string                      `gorm:"index:idx_queue,priority:1"`
	Message     []byte                      `gorm:"type:bytea"`
	AckTime     *time.Time                  `gorm:"type:timestamp; index:idx_ack_time,priority:1"`
	PublishTime time.Time                   `gorm:"type:timestamp; index:idx_queue, sort:asc, priority:2"`
	Errors      shared.ArrayJsonB[ErrorAgg] `gorm:"type:jsonb"`
}
