// Package job contains background job definitions for worker processes.
package job

type CronJob interface {
	Name() string
	Cron() string
	Run() error
}
