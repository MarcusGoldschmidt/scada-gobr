package providers

import (
	"context"
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

var DefaultTimeProvider TimeProvider = UtcTimeProvider{}

func SetTimeProvider(provider TimeProvider) {
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

type LocationTimeProvider struct {
	location *time.Location
}

func NewLocationTimeProvider(location *time.Location) *LocationTimeProvider {
	return &LocationTimeProvider{location: location}
}

func (p LocationTimeProvider) GetCurrentTime() time.Time {
	return time.Now().In(p.location)
}

type UtcTimeProvider struct {
}

func (p UtcTimeProvider) GetCurrentTime() time.Time {
	return time.Now().In(time.UTC)
}

func TimeProviderFromTimeZone(zone *time.Location) TimeProvider {
	if zone == nil {
		return TimeProvider(&UtcTimeProvider{})
	}

	return NewLocationTimeProvider(zone)
}
