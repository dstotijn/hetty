package bolt

import (
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

// Interface guard.
var _ bolt.Logger = (*Logger)(nil)

type Logger struct {
	*zap.SugaredLogger
}

// Warning implements bbolt.Logger.
func (l *Logger) Warning(v ...interface{}) {
	l.Warn(v...)
}

// Warningf implements bbolt.Logger.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}
