package postgres

import (
	"context"
	"encoding/json"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/queue"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SqlPostgresJobQueue struct {
	db           *gorm.DB
	timeProvider providers.TimeProvider
}

type queueMessage struct {
	id   string
	ctx  context.Context
	data any
}

func (q queueMessage) Id() string {
	return q.id
}

func (q queueMessage) Ctx() context.Context {
	return q.ctx
}

func (q queueMessage) Data() any {
	return q.data
}

func NewSqlPostgresJobQueue(db *gorm.DB, timeProvider providers.TimeProvider) *SqlPostgresJobQueue {
	return &SqlPostgresJobQueue{
		db:           db,
		timeProvider: timeProvider,
	}
}

func (s *SqlPostgresJobQueue) Setup() error {
	err := s.db.Migrator().AutoMigrate(&QueueMessages{})
	if err != nil {
		return err
	}

	return nil
}

func (s *SqlPostgresJobQueue) Enqueue(ctx context.Context, queue string, data any) error {
	db := s.db.WithContext(ctx)

	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"TraceParent": util.TraceParent(ctx),
	}

	return db.Create(&QueueMessages{
		Id:          uuid.New(),
		QueueName:   queue,
		Message:     marshal,
		AckTime:     nil,
		Errors:      make([]ErrorAgg, 0),
		PublishTime: s.timeProvider.GetCurrentTime(),
		Headers:     headers,
	}).Error
}

func (s *SqlPostgresJobQueue) Dequeue(ctx context.Context, queueName string, length uint) ([]queue.Message, error) {
	db := s.db.WithContext(ctx)

	currentTime := s.timeProvider.GetCurrentTime()

	var queryResponse []*QueueMessages

	err := db.Model(&QueueMessages{}).Where("queue_name = ? AND publish_time < ? AND ack_time IS NULL", queueName, currentTime).
		Order("publish_time ASC").
		Limit(int(length)).
		Scan(&queryResponse).Error

	if err != nil {
		return nil, err
	}

	responseList := make([]queue.Message, len(queryResponse))

	for i, msg := range queryResponse {
		ctx, err = util.FromTraceParent(context.Background(), msg.Headers["TraceParent"])
		if err != nil {
			ctx = context.Background()
		}

		responseList[i] = queueMessage{
			id:   msg.Id.String(),
			ctx:  ctx,
			data: msg,
		}
	}

	return responseList, nil
}

func (s *SqlPostgresJobQueue) Ack(ctx context.Context, queueName string, messageId string) {
	db := s.db.WithContext(ctx)

	var data *QueueMessages
	result := db.Model(&data).Where("id = ? and queue_name = ?", messageId, queueName).First(&data)

	if result.Error != nil {
		return
	}

	now := s.timeProvider.GetCurrentTime()

	data.AckTime = &now

	db.Save(data)
}

func (s *SqlPostgresJobQueue) Nack(ctx context.Context, queueName string, messageId string, err error) {
	db := s.db.WithContext(ctx)

	var data *QueueMessages
	result := db.Model(&data).Where("id = ? and queue_name = ?", messageId, queueName).First(&data)

	if result.Error != nil {
		return
	}

	data.Errors = append(data.Errors, ErrorAgg{
		Error:        err.Error(),
		CreationTime: s.timeProvider.GetCurrentTime(),
	})

	db.Save(data)
}
