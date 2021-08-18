package uinit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

// Modules registry must be configured by the running application
var Modules = map[string]Module{}

// A Task is an individual Script item
type Task struct {
	Name   string
	Module string
	Args   yaml.Node
	Loop   []string
}

// templateNode this is recursive
func (t *Task) templateNode(ctx *ModuleContext, a *yaml.Node) (err error) {
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
		t.templateNode(ctx, n)
	}
	return
}

// Run an individual task
func (t *Task) Run(ctx *ModuleContext) (err error) {
	var mod Module
	var ok bool
	if mod, ok = Modules[t.Module]; !ok {
		return fmt.Errorf("no module by the name of: %s", t.Module)
	}

	if err = t.templateNode(ctx, &t.Args); err != nil {
		return fmt.Errorf("failed templating arguments: %v", err)
	}

	args := mod.Args()
	if err = t.Args.Decode(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %v", err)
	}
	pre := ctx.Log.Prefix()
	ctx.Log.SetPrefix(pre + t.Module + ": ")
	err = mod.Run(ctx, args)
	ctx.Log.SetPrefix(pre)
	return
}

// A Script object represents a runnable uinit script
type Script struct {
	Context *ModuleContext
	Tasks   []Task
}

func NewScript(data []byte, logger *log.Logger) (s *Script, err error) {
	if logger == nil {
		logger = log.Default()
	}
	s = &Script{
		Context: &ModuleContext{
			Log:  logger,
			Vars: NewSimpleKV(),
		},
		Tasks: []Task{},
	}
	if err = yaml.Unmarshal(data, &s.Tasks); err != nil {
		return nil, fmt.Errorf("failed to parse script: %v", err)
	}
	return
}

func NewScriptFromFile(file string, logger *log.Logger) (s *Script, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return nil, fmt.Errorf("failed to read script file: %v", err)
	}
	return NewScript(data, logger)
}

// Validate test whether a script is sane and ready to run
func (s *Script) Validate() (err error) {
	if s.Tasks == nil {
		return fmt.Errorf("uninitialized task list")
	}
	if s.Context == nil {
		return fmt.Errorf("nil context")
	}
	if s.Context.Log == nil {
		return fmt.Errorf("no logger provided")
	}
	if s.Context.Vars == nil {
		return fmt.Errorf("no KeyValue provided")
	}
	return
}

func (s *Script) Run() (err error) {
	if err = s.Validate(); err != nil {
		return err
	}
	succeed := 0
	total := len(s.Tasks)
	s.Context.Log.Printf("starting uinit script with %d tasks...", total)
	for i, t := range s.Tasks {
		if t.Loop == nil {
			t.Loop = []string{""}
		}
		s.Context.Log.Printf("(%d/%d) task: %s: %s", i+1, total, t.Module, t.Name)
		for _, item := range t.Loop {
			if item != "" {
				s.Context.Log.Printf("(item: %s)", item)
			}
			s.Context.Vars.Set("item", item)
			// we have to copy, otherwise we overwrite our template
			args := &yaml.Node{}
			copier.CopyWithOption(args, &t.Args, copier.Option{DeepCopy: true})
			if err = t.Run(s.Context); err != nil {
				s.Context.Log.Printf("(%d/%d) task failed: %v", i+1, total, err)
			} else {
				s.Context.Log.Printf("(%d/%d) task succeed.", i+1, total)
				succeed++
			}
		}
	}
	s.Context.Log.Printf("(%d/%d) tasks succeeded.  Script finished.  Exiting.", succeed, len(s.Tasks))
	return
}
