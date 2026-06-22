package cronjob

import "go.uber.org/zap"

type logger struct {
	l *zap.SugaredLogger
}

func newLogger(zapLogger *zap.Logger) *logger {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	return &logger{
		l: zapLogger.With(zap.String("component", "cronjob")).Sugar(),
	}
}

func (l *logger) Debug(msg string, args ...any) {
	l.l.Debugf(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.l.Infof(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.l.Warnf(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.l.Errorf(msg, args...)
}
