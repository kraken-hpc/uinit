package setvar

import (
	"fmt"

	"github.com/jlowellwofford/uinit"
)

var _ uinit.Module = (*SetVar)(nil)

// SetVar statically sets a variable
type SetVar struct{}

// Args describes yaml arguments we take
type Args struct {
	Key   string
	Value string
}

// Run the module
func (*SetVar) Run(ctx *uinit.ModuleContext, iargs interface{}) (err error) {
	args, ok := iargs.(*Args)
	if !ok {
		return fmt.Errorf("invalid argument structure")
	}
	ctx.Vars.Set(args.Key, args.Value)
	return
}

// Args returns a struct pointer describing our module's argument structure
func (*SetVar) Args() interface{} {
	return &Args{}
}
