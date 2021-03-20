package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/kraken-hpc/uinit"

	"gopkg.in/yaml.v3"
)

var config = struct {
	scriptFile string
	logFile    string
}{"uinit.script", "uinit.log"}

func usage() {
	fmt.Printf(`
Usage: uinit [<script>]
	default scriptfile: %s
`, config.scriptFile)
}

// a task is an individual item in the yaml script
type task struct {
	Name   string
	Module string
	Args   yaml.Node
}

// templateNode this is recursive
func templateNode(ctx *uinit.ModuleContext, a *yaml.Node) (err error) {
	// too hacky depend on this string never changing?
	if a.ShortTag() == "!!str" {
		c := a.Content
		t := template.Must(template.New("vars").Parse(a.Value)) // we could probably find a clever way to reuse this
		w := bytes.NewBuffer([]byte{})
		if err = t.Execute(w, ctx.Vars.GetMap()); err != nil {
			return fmt.Errorf("failed to template value: %s: %v", a.Value, err)
		}
		if err = a.Encode(w.String()); err != nil {
			return fmt.Errorf("failed to (re)encode value: %v", err)
		}
		a.Content = c
	}
	for _, n := range a.Content {
		templateNode(ctx, n)
	}
	return
}

func run(ctx *uinit.ModuleContext, m string, a *yaml.Node) (err error) {
	var mod uinit.Module
	var ok bool
	if mod, ok = modules[m]; !ok {
		return fmt.Errorf("no module by the name of: %s", m)
	}

	if err = templateNode(ctx, a); err != nil {
		return fmt.Errorf("failed templating arguments: %v", err)
	}

	args := mod.Args()
	if err = a.Decode(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %v", err)
	}
	pre := ctx.Log.Prefix()
	ctx.Log.SetPrefix(pre + m + ": ")
	err = mod.Run(ctx, args)
	ctx.Log.SetPrefix(pre)
	return
}

func main() {
	if len(os.Args) > 2 {
		usage()
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		config.scriptFile = os.Args[1]
	}

	log.SetPrefix("uinit: ")
	log.SetFlags(log.Lmsgprefix | log.Ltime | log.Ldate)
	log.Printf("using script at: %s", config.scriptFile)
	sout, err := ioutil.ReadFile(config.scriptFile)
	if err != nil {
		log.Fatalf("failed to read script file: %v", err)
	}

	script := []task{}
	if err = yaml.Unmarshal(sout, &script); err != nil {
		log.Fatalf("failed to parse script: %v", err)
	}

	log.Printf("using log file: %s", config.logFile)
	logFile, err := os.OpenFile(config.logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	ctx := &uinit.ModuleContext{
		Vars: uinit.NewSimpleKV(),
		Log:  log.New(logFile, "uinit: ", log.Lmsgprefix|log.Ltime|log.Ldate),
	}

	succeed := 0
	total := len(script)
	ctx.Log.Printf("starting uinit script with %d tasks...", total)
	for i, t := range script {
		ctx.Log.Printf("(%d/%d) task: %s: %s", i+1, total, t.Module, t.Name)
		if err = run(ctx, t.Module, &t.Args); err != nil {
			ctx.Log.Printf("(%d/%d) task failed: %v", i+1, total, err)
		} else {
			ctx.Log.Printf("(%d/%d) task succeed.", i+1, total)
			succeed++
		}
	}
	ctx.Log.Printf("(%d/%d) tasks succeeded.  Script finished.  Exiting.", succeed, len(script))
}
