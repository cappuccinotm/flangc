package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"
	"github.com/cappuccinotm/flangc/app/cmd"
)

// Options describes command line options for an application.
type Options struct {
	Run   cmd.Run `command:"run"`
	Debug bool    `long:"dbg" env:"DEBUG" description:"turn on debug mode"`
}

var version = "unknown"

func main() {
	fmt.Printf("flangc, version: %s\n", version)

	var opts Options

	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(cmd flags.Commander, args []string) error {
		setupLog(opts.Debug)

		if err := cmd.Execute(args); err != nil {
			log.Printf("[ERROR] failed to execute command %+v", err)
		}
		return nil
	}

	// after failure command does not return non-zero code
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println("failed to run command: %v", err)
			os.Exit(1)
		}
	}
}

func setupLog(dbg bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "INFO",
		Writer:   os.Stdout,
	}

	logFlags := log.Ldate | log.Ltime

	if dbg {
		logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
		filter.MinLevel = "DEBUG"
	}

	log.SetFlags(logFlags)
	log.SetOutput(filter)
}
