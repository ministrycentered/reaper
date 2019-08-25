package reaper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// CallbackHandler is the function signature for app callbacks
type CallbackHandler func(*Context) error

// App contains the configuration for a command line application
type App struct {
	name        string
	output      *log.Logger
	errLogger   *log.Logger
	commands    map[string]*Command
	before      []callback
	Description string
	Version     string
}

type callback struct {
	handler  CallbackHandler
	excluded []string
}

func (c callback) isExcluded(name string) bool {
	for _, key := range c.excluded {
		if key == name {
			return true
		}
	}
	return false
}

// NewApp returns a new app ready to be configured
func NewApp(name string) *App {
	app := &App{
		name:      name,
		output:    log.New(os.Stdout, "", 0),
		errLogger: log.New(os.Stderr, "❗️ ", 0),
		commands:  make(map[string]*Command, 0),
		before:    make([]callback, 0),
		Version:   "1.0.0",
	}

	app.Command("help", helpCommandHandler).Configure(func(c *Command) {
		c.isInternal = true
		c.Description = "prints this help dialogue"
	})

	app.Command("version", versionHandler).Configure(func(c *Command) {
		c.isInternal = true
		c.Description = "print the current version"
	})

	app.Command("_commands", outputCommandsHandler).Configure(func(c *Command) {
		c.isInternal = true
		c.Private = true
	})

	return app
}

// Command creates a new command with the handler and adds it to the app.
func (a *App) Command(name string, handler CommandHandler) *Command {
	cmd := newCommand(name, handler)
	a.commands[name] = cmd
	return cmd
}

// Before adds a function handler
//
// This auto-excludes the help & version commands
func (a *App) Before(fn CallbackHandler) {
	a.BeforeExcluding(fn, "help", "version")
}

// BeforeExcluding adds a function handler but doesn't run it for the passed command names
func (a *App) BeforeExcluding(fn CallbackHandler, exclude ...string) {
	a.before = append(a.before, callback{handler: fn, excluded: exclude})
}

// Execute takes a list of arguments and runs the command that matches.
func (a *App) Execute(args []string) error {
	if len(args) < 1 {
		return errors.New("Too few arguments, must have at least one. Try `help`")
	}
	name := args[0]
	args = args[1:]

	cmd, ok := a.commands[name]
	if !ok {
		return fmt.Errorf("no command with the name %s. Try `help`", name)
	}

	ctx := newContext(a, name, args)
	flags, err := cmd.flagSet(ctx, name)
	if err != nil {
		return err
	}

	err = flags.Parse(args)
	if err != nil {
		return err
	}

	cmd.handleArgs(ctx, flags.Args())

	for _, callback := range a.before {
		if callback.isExcluded(name) {
			continue
		}
		err = callback.handler(ctx)
		if err != nil {
			return err
		}
	}

	return cmd.handler(ctx)
}

func (a *App) commandNames() []string {
	names := []string{}
	for name, cmd := range a.commands {
		if cmd.Private {
			continue
		}
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func versionHandler(c *Context) error {
	c.Output(c.app.Version)
	return nil
}

func outputCommandsHandler(c *Context) error {
	c.Output(strings.Join(c.app.commandNames(), " "))

	return nil
}
