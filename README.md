# uinit

Uinit is a simple golang init process that is intended to be run from within u-root (but can run elsewhere too, say, as a light-weight container init).

## Overview

Uinit is built to be *scriptable* with a simple YAML-format file.  For an example script, see: [cmds/uinit/uinit.script](cmds/uinit/uinit.script).

By default, uinit reads the script at the relative path `./uinit.script` and writes tot he log file `./uinit.log`. An alternative script can be specified on the commandline: `./uinit <path_to_script>`

Uinit maintains a *key/value store* that allows for storing and recalling variables within the script.  Variables can be accessed by placing strings of the format `"{{.<key>}}"` anywhere in the args section of a task.  Technically, anything that can be processed by Go's `text/template` package can be placed here.

Uinit reads the script as a sequence of tasks.  Each task calls a "module", and passes "args" to the module.  Currently, we have the following modules:

- echo : Echo something to the log
- command : Execute a command.  This can execute in foreground, background, or call Exec (which terminates uinit).  Commands can be executed within a shell or directly.
- setvar : Sets a variable in the key/value store that can be later referenced.
- cmdline : Reads a file containing commandline arguments of the style: `uinit.<key>=<val>` adn sets `<key> = <val>` in the key/value store.  If formated as `uinit.key`, then `val="true"` implicitly.  This is intended primarily to read variables from `/proc/cmdline` which is the default file to parse, allowing for uinit configuration by kernel parameters.  The prefix, which defaults to `uinit` can be set to anything.
- mount : Mount a filesystem.
- script : Run a subscript, either directly specified or in another script file.

More modules to come! (maybe?)

For more details on modules, see the example scripts, or look for `README.md` files in the `modules/<module>` directory.

## Building/using

1. `go get -u github.com/kraken-hpc/uinit/cmds/uinit`
   This will build `uinit` and place it in `$GOPATH/bin`.
2. Write a script (see example `*.script` files)
3. Run it, and either name the script `./uinit.script`, or pass it as a commandline argument.

## Using as a package

As of v0.2.0, `uinit` can also be used as a Go package.  This package will let your application run `uinit` scripts natively.

As an example, to load and run a script from a file:

```go
s, err := uinit.NewScriptFromFile("myscript.script", logger)
if err != nil {
   log.Fatalf("failed to load script: %v", err)
}
if err = s.Run(); err != nil {
   log.Fatalf("failed to run script: %v", err)
}
```

For more details, see the package documentation at [https://pkg.go.dev/github.com/kraken-hpc/uinit](https://pkg.go.dev/github.com/kraken-hpc/uinit).


## Building into u-root

You can build uinit directly into [u-root](https://github.com/u-root/u-root) by adding the following commandline options to your u-root build:
`-uinitcmd /bbin/uinit github.com/kraken-hpc/uinit/cmds/uinit`
Then make sure that you have a `uinit.script` in the root of your `base` cpio.

If you intend to boot something with `systemd`out of this, you'll want to:
1. have `uroot.systemd` on your kernel commandline
2. put a built `uinit` named `inito` in the root of your `base cpio`

## Enjoy
