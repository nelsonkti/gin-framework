package cron

type Task interface {
	Rule() string
	Run()
	Name() string
}
