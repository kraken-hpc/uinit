package cmdline

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/jlowellwofford/uinit"
)

var _ uinit.Module = (*CmdLine)(nil)

// CmdLine parses variables out of a commandline indicated by a filename
// This is mostly intended to get variables from boot-time /proc/cmdline
type CmdLine struct{}

const (
	defaultPrefix   = "uinit" // variables named uinit.<var>
	defaultFilename = "/proc/cmdline"
)

// Args describes yaml arguments we take
type Args struct {
	Filename string `yaml:"filename,omitempty"`
	Prefix   string `yaml:"prefix,omitempty"`
}

// Run the module
func (*CmdLine) Run(ctx *uinit.ModuleContext, iargs interface{}) (err error) {
	args, ok := iargs.(*Args)
	if !ok {
		return fmt.Errorf("invalid argument structure")
	}
	if args.Filename == "" {
		args.Filename = defaultFilename
	}
	if args.Prefix == "" {
		args.Prefix = defaultPrefix
	}
	var re *regexp.Regexp
	/*
	 * uinit.key=val -> key = val
	 * uinit.key -> key = "true"
	 * uinit.key= -> key = ""
	 */
	if re, err = regexp.Compile(`^` + args.Prefix + `\.(\w+)(=(.*))?$`); err != nil {
		return fmt.Errorf("invalid prefix: %v", err)
	}
	var text []byte
	if text, err = ioutil.ReadFile(args.Filename); err != nil {
		return fmt.Errorf("unable to read cmdline file: %v", err)
	}
	for _, o := range uinit.SplitCommandLine(string(text)) {
		m := re.FindAllStringSubmatch(o, 4)
		if len(m) > 0 {
			var k, v string
			k = m[0][1]
			v = m[0][3]
			// we have a match
			if m[0][3] == "" && m[0][2] == "" {
				// this is an implicit true
				v = "true"
			}
			ctx.Log.Printf("setting: %s = %s", k, v)
			ctx.Vars.Set(k, v)
		}
	}
	return
}

// Args returns a struct pointer describing our module's argument structure
func (*CmdLine) Args() interface{} {
	return &Args{}
}
