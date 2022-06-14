package providers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

const TimeProviderCtxKey = "timeProvider"

func GetTimeProviderFromCtx(ctx context.Context) TimeProvider {
	if ctx == nil {
		return TimeProvider(&UtcTimeProvider{})
	}

	if provider, ok := ctx.Value(TimeProviderCtxKey).(TimeProvider); ok {
		return provider
	}
	return TimeProvider(&UtcTimeProvider{})
}

var DefaultTimeProvider = UtcTimeProvider{}

func SetTimeProvider(provider UtcTimeProvider) {
	DefaultTimeProvider = provider
}

type TimeProvider interface {
	GetCurrentTime() time.Time
}

type SimpleTimeProvider struct {
}

func (p SimpleTimeProvider) GetCurrentTime() time.Time {
	return time.Now()
}

type LocationTimeProvider string

func (p LocationTimeProvider) GetCurrentTime() time.Time {
	location, err := time.LoadLocation(string(p))
	if err != nil {
		log.Errorf("%s", err.Error())
		return time.Now()
	}
	return time.Now().In(location)
}

type UtcTimeProvider struct {
}

func (p UtcTimeProvider) GetCurrentTime() time.Time {
	return time.Now().In(time.UTC)
}

func TimeProviderFromTimeZone(zone string) TimeProvider {
	if zone == "" {
		return TimeProvider(&UtcTimeProvider{})
	}

	return LocationTimeProvider(zone)
}
