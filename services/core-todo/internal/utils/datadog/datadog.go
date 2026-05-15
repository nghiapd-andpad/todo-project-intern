package datadog

import (
	"net"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type datadogZapLogger struct {
	zapLogger *zap.Logger
}

func newDatadogZapLogger(zapLogger *zap.Logger) tracer.Logger {
	return &datadogZapLogger{zapLogger: zapLogger}
}

func (l *datadogZapLogger) Log(msg string) {
	l.zapLogger.Info(msg)
}

func InitDatadog(
	ddApmEnabled bool,
	ddAgentHost string,
	ddAgentPort string,
	ddDogstatsdPort string,
	zapLogger *zap.Logger) (func(), error) {
	if !ddApmEnabled {
		return func() {}, nil
	}
	if zapLogger == nil {
		return func() {}, entity.NewInternal("failed to get zapLogger")
	}
	if ddAgentPort == "" || ddAgentHost == "" || ddDogstatsdPort == "" {
		return func() {}, entity.NewInternal(
			"failed to get params").
			WithDetail("host", ddAgentHost).
			WithDetail("port", ddAgentPort)
	}

	addr := net.JoinHostPort(ddAgentHost, ddAgentPort)
	dogstatsdAddr := net.JoinHostPort(ddAgentHost, ddDogstatsdPort)
	logger := newDatadogZapLogger(zapLogger)

	if err := tracer.Start(
		tracer.WithAgentAddr(addr),
		tracer.WithDogstatsdAddr(dogstatsdAddr),
		tracer.WithLogger(logger),
		tracer.WithRuntimeMetrics(),
	); err != nil {
		return func() {}, entity.NewInternal("failed to start tracer").WithDetail("desc", err.Error())
	}

	return tracer.Stop, nil
}
