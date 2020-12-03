package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unicode"

	"github.com/jlowellwofford/uinit"
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
		cmdArgv = split(args.Cmd)
		if cmdPath, err = exec.LookPath(cmdArgv[0]); err != nil {
			return fmt.Errorf("command not found: %s", cmdArgv[0])
		}
	}

	if args.Exec {
		return syscall.Exec(cmdPath, cmdArgv[1:], os.Environ())
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

//Split strings on spaces except when a space is within a quoted, bracketed, or braced string.
//Supports nesting multiple brackets or braces.
func split(s string) []string {
	lastRune := map[rune]int{}
	f := func(c rune) bool {
		switch {
		case lastRune[c] > 0:
			lastRune[c]--
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastRune[c]++
			return false
		case c == '[':
			lastRune[']']++
			return false
		case c == '{':
			lastRune['}']++
			return false
		case mapGreaterThan(lastRune, 0):
			return false
		default:
			return c == ' '
		}
	}
	return strings.FieldsFunc(s, f)
}

// mapGreaterThan ranges across the provided map[rune]int looking for any values greater than
func mapGreaterThan(runes map[rune]int, g int) bool {
	for _, i := range runes {
		if i > g {
			return true
		}
	}
	return false
}
