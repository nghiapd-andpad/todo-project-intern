// Package gateway defines outbound ports used by application services.
package gateway

import "context"

type Scheduler interface {
	ScheduleCron(ctx context.Context, name string, cronExpr string, task any, parameters ...any) error
	Start()
	Stop() error
}
