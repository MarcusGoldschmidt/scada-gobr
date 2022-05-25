package logger

import (
	log "github.com/sirupsen/logrus"
	"io"
)

// SimpleLogger ...
type SimpleLogger struct {
	Logger   *log.Logger
	LogLevel LogLevel
}

// NewSimpleLogger ...
func NewSimpleLogger(name string, out io.Writer) Logger {
	return &SimpleLogger{
		Logger:   log.New(),
		LogLevel: logLevelFromEnvironment(),
	}
}

// NewSimpleLoggerWithLevel ...
func NewSimpleLoggerWithLevel(name string, out io.Writer, level LogLevel) Logger {
	return &SimpleLogger{
		Logger:   log.New(),
		LogLevel: level,
	}
}

// CloneWithLevel ...
func (l *SimpleLogger) CloneWithLevel(level LogLevel) Logger {
	return &SimpleLogger{
		Logger:   l.Logger,
		LogLevel: level,
	}
}

// Errorf ...
func (l *SimpleLogger) Errorf(f string, v ...interface{}) {
	if l.LogLevel <= LogError {
		l.Logger.Printf("ERROR: "+f, v...)
	}
}

// Warningf ...
func (l *SimpleLogger) Warningf(f string, v ...interface{}) {
	if l.LogLevel <= LogWarn {
		l.Logger.Printf("WARNING: "+f, v...)
	}
}

// Infof ...
func (l *SimpleLogger) Infof(f string, v ...interface{}) {
	if l.LogLevel <= LogInfo {
		l.Logger.Printf("INFO: "+f, v...)
	}
}

// Debugf ...
func (l *SimpleLogger) Debugf(f string, v ...interface{}) {
	if l.LogLevel <= LogDebug {
		l.Logger.Printf("DEBUG: "+f, v...)
	}
}
