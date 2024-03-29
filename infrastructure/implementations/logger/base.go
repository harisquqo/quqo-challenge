package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/logger_repository"
	honeycomb "github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger/honeycomb_implementation"
	zap "github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger/zap_implementation"
	"go.opentelemetry.io/otel/trace"
)

// LoggerRepo is a logger repository that can use multiple channels
type LoggerRepo struct {
	loggers []logger_repository.LoggerRepository
	c *gin.Context
	honeycombRepo *honeycomb.HoneycombRepo
}

const (
	Zap = "Zap"
	Honeycomb = "Honeycomb"
)

type Option func(*LoggerRepo)

// NewLoggerRepository creates a new logger repository based on the specified channels
func NewLoggerRepository(channels []string) (LoggerRepo) {
	var loggers []logger_repository.LoggerRepository
	var honeycombRepo *honeycomb.HoneycombRepo
	for _, channel := range channels {
		switch channel {
		case Zap:
			loggers = append(loggers, zap.NewZapRepository())
		case Honeycomb:
			honeycombRepo = honeycomb.NewHoneycombRepository()
			loggers = append(loggers, honeycombRepo)
		default:
			// You might want to log or handle unsupported channels
			continue
		}
	}
	
	return LoggerRepo{loggers: loggers, c: nil, honeycombRepo: honeycombRepo}
}

// Debug logs a debug message
func (l *LoggerRepo) Debug(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Debug(msg, fields)
	}
}

// Info logs an info message
func (l *LoggerRepo) Info(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Info(msg, fields)
	}
}

// Warn logs a warning message
func (l *LoggerRepo) Warn(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Warn(msg, fields)
	}
}

// Error logs an error message
func (l *LoggerRepo) Error(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Error(msg, fields)
	}
}

// Fatal logs a fatal message
func (l *LoggerRepo) Fatal(msg string, fields map[string]interface{}) {
	for _, logger := range l.loggers {
		logger.Fatal(msg, fields)
	}
}

// Start function 
func (l *LoggerRepo) Start(c *gin.Context, info string, options ...Option) trace.Span {
	l.c = c

	var span trace.Span
	for _, logger := range l.loggers {
		_, ok := logger.(*honeycomb.HoneycombRepo)
		if ok {
			span = logger.Start(c, info)
		} else {
			logger.Start(c, info)
		}
	}

	// Execute optional functions
    for _, opts := range options {
        opts(l)
	}

	return span
}

func (l *LoggerRepo) SetContextWithSpanFunc() Option {
    return func(l *LoggerRepo) {
        l.c.Set("otel_context", trace.ContextWithSpan(l.c, l.honeycombRepo.GetSpan()))
    }
}

func (l *LoggerRepo) SetContextWithSpan(span trace.Span) {
	l.c.Set("otel_context", trace.ContextWithSpan(l.c, span))
}