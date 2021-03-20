package command

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/kraken-hpc/uinit"
)

var _ uinit.Module = (*Command)(nil)

var shell = []string{"/bin/sh", "-c"}

// Command executes a command
type Command struct{}

// Args describes yaml arguments we take
type Args struct {
	Cmd        string
	Background bool
	Exec       bool
	Shell      bool
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

	var cmdPath string
	var cmdArgv []string

	if args.Shell {
		// ok, set us up to run in shell
		cmdArgv = append(shell, args.Cmd)
		cmdPath = cmdArgv[0]
	} else {
		cmdArgv = uinit.SplitCommandLine(args.Cmd)
		if cmdPath, err = exec.LookPath(cmdArgv[0]); err != nil {
			return fmt.Errorf("command not found: %s", cmdArgv[0])
		}
	}

	if args.Exec {
		return syscall.Exec(cmdPath, cmdArgv, os.Environ())
	}

	c := exec.Command(cmdPath, cmdArgv[1:]...)
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
