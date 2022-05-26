package providers

import (
	log "github.com/sirupsen/logrus"
	"time"
)

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
