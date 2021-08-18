package script

import (
	"fmt"

	"github.com/kraken-hpc/uinit"
)

var _ uinit.Module = (*Script)(nil)

// Script simply echoes some text
type Script struct{}

// Args describes yaml arguments we take
type Args struct {
	File  string
	Tasks []uinit.Task
}

// Run the module
func (*Script) Run(ctx *uinit.ModuleContext, iargs interface{}) (err error) {
	args, ok := iargs.(*Args)
	if !ok {
		return fmt.Errorf("invalid argument structure")
	}
	if args.Tasks == nil {
		args.Tasks = []uinit.Task{}
	}
	var s *uinit.Script
	if args.File != "" {
		if s, err = uinit.NewScriptFromFile(args.File, ctx.Log); err != nil {
			return fmt.Errorf("failed to read script file: %v", err)
		}
	} else {
		if s, err = uinit.NewScript([]byte{}, ctx.Log); err != nil {
			return fmt.Errorf("failed to create script: %v", err)
		}
	}
	// arbitrarily, specified tasks go before file tasks
	s.Tasks = append(args.Tasks, s.Tasks...)
	return s.Run()
}

// Args returns a struct pointer describing our module's argument structure
func (*Script) Args() interface{} {
	return &Args{}
}
