package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kraken-hpc/uinit"
	"github.com/kraken-hpc/uinit/modules"
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

func main() {
	if len(os.Args) > 2 {
		usage()
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		config.scriptFile = os.Args[1]
	}
	modules.InitAll()

	log.SetPrefix("uinit: ")
	log.SetFlags(log.Lmsgprefix | log.Ltime | log.Ldate)
	log.Printf("using script at: %s", config.scriptFile)
	script, err := uinit.NewScriptFromFile(config.scriptFile, nil)
	if err != nil {
		log.Fatalf("failed to read script: %v", err)
	}

	log.Printf("using log file: %s", config.logFile)
	logFile, err := os.OpenFile(config.logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	script.Context.Log = log.New(io.MultiWriter(os.Stdout, logFile), "uinit: ", log.Lmsgprefix|log.Ltime|log.Ldate)
	if err = script.Run(); err != nil {
		log.Fatalf("script failed to run: %v", err)
	}
}
