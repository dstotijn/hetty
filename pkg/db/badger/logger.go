package badger

import (
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// Interface guard.
var _ badger.Logger = (*Logger)(nil)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(l *zap.SugaredLogger) *Logger {
	return &Logger{l}
}

func (l *Logger) Warningf(template string, args ...interface{}) {
	l.Warnf(template, args)
}
