package command

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/jlowellwofford/uinit"
)

var _ uinit.Module = (*Command)(nil)

var shell = "/bin/sh"
var shellopts = []string{"-c"}

// Command executes a command
type Command struct{}

// Args describes yaml arguments we take
type Args struct {
	Cmd        string
	Background bool
	Exec       bool
	Argv       []string `yaml:"argv,omitempty"`
}

// Run the module
func (*Command) Run(ctx *uinit.ModuleContext, iargs interface{}) (err error) {
	args, ok := iargs.(*Args)
	if !ok {
		return fmt.Errorf("invalid argument structure")
	}

	if args.Cmd == "" {
		// noop
		return fmt.Errorf("no cmd specified")
	}
	var path string
	if path, err = exec.LookPath(args.Cmd); err != nil {
		// maybe we should run in a shell? Not if argv is specified
		if len(args.Argv) != 0 {
			return fmt.Errorf("command not found: %s", args.Cmd)
		}
		// ok, set us up to run in shell
		args.Argv = append(shellopts, args.Cmd)
		args.Cmd = shell
	} else {
		args.Cmd = path
	}

	if args.Exec {
		return syscall.Exec(args.Cmd, args.Argv, os.Environ())
	}

	c := exec.Command(args.Cmd, args.Argv...)
	if args.Background {
		c.Stdout = ctx.Log.Writer()
		c.Stderr = ctx.Log.Writer()
		c.Stdin = nil
		return c.Start()
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// Args returns a struct pointer describing our module's argument structure
func (*Command) Args() interface{} {
	return &Args{}
}
