package mount

import (
	"fmt"

	"github.com/bensallen/rbd/pkg/mount"
	"github.com/kraken-hpc/uinit"
)

var _ uinit.Module = (*Mount)(nil)

// Mount mounts a filesystem
type Mount struct{}

// Args describes yaml arguments we take
type Args struct {
	Src  string
	Dest string
	Type string
	Opts []string
}

// Run the module
func (*Mount) Run(ctx *uinit.ModuleContext, iargs interface{}) (err error) {
	args, ok := iargs.(*Args)
	if !ok {
		return fmt.Errorf("invalid argument structure")
	}
	if err = mount.Mount(args.Src, args.Dest, args.Type, args.Opts); err != nil {
		return fmt.Errorf("mount failed: %v", err)
	}
	ctx.Log.Printf("mounted %s -> %s (type: %s)", args.Src, args.Dest, args.Type)
	return
}

// Args returns a struct pointer describing our module's argument structure
func (*Mount) Args() interface{} {
	return &Args{}
}
