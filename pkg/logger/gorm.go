package logger

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	*log.Logger

	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func NewGormLogger() *GormLogger {
	logger := log.New()

	logger.Level = log.InfoLevel

	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	return &GormLogger{
		Logger:                logger,
		SkipErrRecordNotFound: false,
	}
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.WithContext(ctx).Infof(s, args)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.WithContext(ctx).Warnf(s, args)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.WithContext(ctx).Errorf(s, args)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := log.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		l.Logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	l.Logger.WithContext(ctx).WithFields(fields).Tracef("%s [%s]", sql, elapsed)
}
