package reaper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

// CallbackHandler is the function signature for app callbacks
type CallbackHandler func(*Context) error

// App contains the configuration for a command line application
type App struct {
	name        string
	output      *log.Logger
	errLogger   *log.Logger
	commands    map[string]*Command
	before      []CallbackHandler
	Description string
	Version     string
}

// NewApp returns a new app ready to be configured
func NewApp(name string) *App {
	app := &App{
		name:      name,
		output:    log.New(os.Stdout, "", 0),
		errLogger: log.New(os.Stderr, "❗️ ", 0),
		commands:  make(map[string]*Command, 0),
		before:    make([]CallbackHandler, 0),
		Version:   "1.0.0",
	}

	app.Command("help", helpCommandHandler).Configure(func(c *Command) {
		c.Description = "prints this help dialogue"
	})

	app.Command("version", versionHandler).Configure(func(c *Command) {
		c.Description = "print the current version"
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
func (a *App) Before(fn CallbackHandler) {
	a.before = append(a.before, fn)
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

	ctx := newContext(a, args)
	flags, err := cmd.flagSet(ctx, name)
	if err != nil {
		return err
	}

	err = flags.Parse(args)
	if err != nil {
		return err
	}

	cmd.handleArgs(ctx, flags.Args())

	for _, handler := range a.before {
		err = handler(ctx)
		if err != nil {
			return err
		}
	}

	return cmd.handler(ctx)
}

func (a *App) commandNames() []string {
	names := []string{}
	for name := range a.commands {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func versionHandler(c *Context) error {
	c.Output(c.app.Version)
	return nil
}
