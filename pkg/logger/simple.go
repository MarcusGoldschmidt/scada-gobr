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
	logger := log.New()

	logger.Out = out

	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	return &SimpleLogger{
		Logger:   logger,
		LogLevel: LogLevelFromEnvironment(),
	}
}

// NewSimpleLoggerWithLevel ...
func NewSimpleLoggerWithLevel(name string, out io.Writer, level LogLevel) Logger {
	return &SimpleLogger{
		Logger:   log.New(),
		LogLevel: level,
	}
}

func (l *SimpleLogger) Tracef(s string, i ...interface{}) {
	if l.LogLevel <= LogError {
		l.Logger.Tracef(s, i...)
	}
}

// Errorf ...
func (l *SimpleLogger) Errorf(f string, v ...interface{}) {
	if l.LogLevel <= LogError {
		l.Logger.Errorf(FileWithLineNum()+" "+f, v...)
	}
}

// Warningf ...
func (l *SimpleLogger) Warningf(f string, v ...interface{}) {
	if l.LogLevel <= LogWarn {
		l.Logger.Warningf(f, v...)
	}
}

// Infof ...
func (l *SimpleLogger) Infof(f string, v ...interface{}) {
	if l.LogLevel <= LogInfo {
		l.Logger.Infof(f, v...)
	}
}

// Debugf ...
func (l *SimpleLogger) Debugf(f string, v ...interface{}) {
	if l.LogLevel <= LogDebug {
		l.Logger.Debugf(f, v...)
	}
}
