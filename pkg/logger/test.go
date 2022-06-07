package logger

import (
	"testing"
)

// SimpleLogger ...
type TestLogger struct {
	Logger   *testing.T
	LogLevel LogLevel
}

func NewTestLogger(t *testing.T) Logger {
	return &TestLogger{
		Logger:   t,
		LogLevel: LogDebug,
	}
}

func (l *TestLogger) Tracef(s string, i ...interface{}) {
	if l.LogLevel <= LogTrace {
		l.Logger.Logf("TRACE: "+s, i...)
	}
}

// Errorf ...
func (l *TestLogger) Errorf(f string, v ...interface{}) {
	if l.LogLevel <= LogError {
		l.Logger.Errorf("ERROR: "+f, v...)
	}
}

// Warningf ...
func (l *TestLogger) Warningf(f string, v ...interface{}) {
	if l.LogLevel <= LogWarn {
		l.Logger.Logf("WARNING: "+f, v...)
	}
}

// Infof ...
func (l *TestLogger) Infof(f string, v ...interface{}) {
	if l.LogLevel <= LogInfo {
		l.Logger.Logf("INFO: "+f, v...)
	}
}

// Debugf ...
func (l *TestLogger) Debugf(f string, v ...interface{}) {
	if l.LogLevel <= LogDebug {
		l.Logger.Logf("DEBUG: "+f, v...)
	}
}
