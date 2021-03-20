package echo

import (
	"fmt"

	"github.com/kraken-hpc/uinit"
)

var _ uinit.Module = (*Echo)(nil)

// Echo simply echoes some text
type Echo struct{}

// Args describes yaml arguments we take
type Args struct {
	Text      string
	Nonewline bool
}

// Run the module
func (*Echo) Run(ctx *uinit.ModuleContext, iargs interface{}) (err error) {
	args, ok := iargs.(*Args)
	if !ok {
		return fmt.Errorf("invalid argument structure")
	}
	ctx.Log.Printf("%s", args.Text)
	return
}

// Args returns a struct pointer describing our module's argument structure
func (*Echo) Args() interface{} {
	return &Args{}
}
