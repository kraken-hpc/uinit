package uinit

import "log"

// KeyValue describes our requirements for a key value store
type KeyValue interface {
	Init()
	Set(string, string) error
	Get(string) (string, error)
	GetDefault(string, string) string
	Clear() error
	GetMap() map[string]string
}

// ModuleContext gets passed to modules when they Run
type ModuleContext struct {
	Vars KeyValue
	Log  *log.Logger
}

// A Module is a named hook that can perform some task in our script
type Module interface {
	Args() interface{}
	Run(ctx *ModuleContext, args interface{}) error
}
