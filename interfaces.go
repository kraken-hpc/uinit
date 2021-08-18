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
// Vars provides a KeyValue store to be used by the script
// Log provides output logging capabilities for the script
type ModuleContext struct {
	Vars KeyValue
	Log  *log.Logger
}

// A Module is a named hook that can perform some task in our script
// A Script consists of a set of Tasks that use Modules to do work
type Module interface {
	Args() interface{}
	Run(ctx *ModuleContext, args interface{}) error
}
