package modules

import (
	"github.com/kraken-hpc/uinit"
	"github.com/kraken-hpc/uinit/modules/cmdline"
	"github.com/kraken-hpc/uinit/modules/command"
	"github.com/kraken-hpc/uinit/modules/echo"
	"github.com/kraken-hpc/uinit/modules/script"
	"github.com/kraken-hpc/uinit/modules/setvar"
)

func InitAll() {
	uinit.Modules = map[string]uinit.Module{
		"echo":    &echo.Echo{},
		"setvar":  &setvar.SetVar{},
		"command": &command.Command{},
		"cmdline": &cmdline.CmdLine{},
		"script":  &script.Script{},
	}
}
