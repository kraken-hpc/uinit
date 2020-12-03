package main

import (
	"github.com/jlowellwofford/uinit"
	"github.com/jlowellwofford/uinit/modules/cmdline"
	"github.com/jlowellwofford/uinit/modules/command"
	"github.com/jlowellwofford/uinit/modules/echo"
	"github.com/jlowellwofford/uinit/modules/setvar"
)

var modules = map[string]uinit.Module{
	"echo":    &echo.Echo{},
	"setvar":  &setvar.SetVar{},
	"command": &command.Command{},
	"cmdline": &cmdline.CmdLine{},
}
