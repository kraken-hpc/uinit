package main

import (
	"github.com/kraken-hpc/uinit"
	"github.com/kraken-hpc/uinit/modules/cmdline"
	"github.com/kraken-hpc/uinit/modules/command"
	"github.com/kraken-hpc/uinit/modules/echo"
	"github.com/kraken-hpc/uinit/modules/mount"
	"github.com/kraken-hpc/uinit/modules/setvar"
)

var modules = map[string]uinit.Module{
	"echo":    &echo.Echo{},
	"setvar":  &setvar.SetVar{},
	"command": &command.Command{},
	"cmdline": &cmdline.CmdLine{},
	"mount":   &mount.Mount{},
}
