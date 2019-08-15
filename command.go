package reaper

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

// CommandHandler is the function that will be called when a user performs a sub-command
type CommandHandler func(*Context) error

// Command represents a single sub-command in the application
type Command struct {
	handler     CommandHandler
	Description string
	flags       map[string]*flagDesc
	args        []*argDesc
	examples    []string
}

type flagDesc struct {
	kind  string
	value interface{}
	usage string
}

type argDesc struct {
	required bool
	name     string
	usage    string
}

func newCommand(name string, handler CommandHandler) *Command {
	return &Command{
		handler:  handler,
		flags:    make(map[string]*flagDesc, 0),
		args:     make([]*argDesc, 0),
		examples: make([]string, 0),
	}
}

// Configure calls the passed function with the command.
//
// This is a convenience method for chaining `app.Command("thing", thingHandler).Configure(func(c *Command) { ... })`
func (c *Command) Configure(fn func(*Command)) *Command {
	fn(c)
	return c
}

// Argument adds a named argument to the command
func (c *Command) Argument(name string, required bool, usage string) {
	c.args = append(c.args, &argDesc{name: name, required: required, usage: usage})
}

// Flag adds a cli flag to the command
func (c *Command) Flag(name, kind string, value interface{}, usage string) {
	c.flags[name] = &flagDesc{kind: kind, value: value, usage: usage}
}

// Example adds and example to the output
func (c *Command) Example(e string) {
	c.examples = append(c.examples, e)
}

func (c *Command) flagSet(ctx *Context, name string) (*flag.FlagSet, error) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for name, desc := range c.flags {
		switch desc.kind {
		case "string":
			ctx.flags[name] = set.String(name, desc.value.(string), desc.usage)
		case "integer":
			ctx.flags[name] = set.Int64(name, desc.value.(int64), desc.usage)
		case "collection":
			collection := &FlagCollection{Values: make([]string, 0)}
			set.Var(collection, name, desc.usage)
			ctx.flags[name] = collection
		default:
			return nil, fmt.Errorf("unknown flag type: %s", desc.kind)
		}
	}

	return set, nil
}

func (c *Command) handleArgs(ctx *Context, args []string) error {
	if len(c.args) < len(args) {
		return fmt.Errorf("expected %d arguments, got %d", len(c.args), len(args))
	}

	for idx, desc := range c.args {
		if desc.required && idx >= len(args) {
			return fmt.Errorf("missing required argument for %s", desc.name)
		}
		if !desc.required && idx >= len(args) {
			break
		}
		ctx.namedArgs[desc.name] = args[idx]
	}

	return nil
}

func (c *Command) flagNames() []string {
	names := []string{}
	for name := range c.flags {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// FlagCollection contains a list of flags
type FlagCollection struct {
	Values []string
}

func (f *FlagCollection) String() string {
	return strings.Join(f.Values, ", ")
}

// Set adds a flag value
func (f *FlagCollection) Set(value string) error {
	f.Values = append(f.Values, value)

	return nil
}
