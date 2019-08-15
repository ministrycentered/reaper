package reaper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Context is used to hold the state of the application
type Context struct {
	CommandName string
	app         *App
	args        []string
	flags       map[string]interface{}
	namedArgs   map[string]string
	storage     map[string]interface{}
}

func newContext(app *App, name string, args []string) *Context {
	return &Context{
		CommandName: name,
		app:         app,
		args:        args,
		flags:       make(map[string]interface{}, 0),
		namedArgs:   make(map[string]string, 0),
		storage:     make(map[string]interface{}, 0),
	}
}

// Flag returns a flag with the passed name. If the flag is not found an error will be returned.
func (c *Context) Flag(name string) (interface{}, error) {
	if val, ok := c.flags[name]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("missing required flag %s", name)
}

// Set stores a value in the context for the key
func (c *Context) Set(key string, value interface{}) {
	c.storage[key] = value
}

// Get returns a stored value from the context
func (c *Context) Get(key string) interface{} {
	return c.storage[key]
}

// FlagString returns a flag string value.
//
// This will panic if the flag does not exist
func (c *Context) FlagString(name string) string {
	value, _ := c.Flag(name)
	return *(value.(*string))
}

func (c *Context) FlagCollection(name string) []string {
	value, _ := c.Flag(name)
	coll := value.(*FlagCollection)
	return coll.Values
}

// Argument returns a named argument
func (c *Context) Argument(name string) string {
	return c.namedArgs[name]
}

// Output writes the string to stdout
func (c *Context) Output(value string) {
	c.app.output.Print(value)
}

// Input asks for user input
func (c *Context) Input(msg string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "%s ", msg)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(text, "\n"), nil
}

// Confirm asks for conformation before continuing
func (c *Context) Confirm(msg string) bool {
	response, _ := c.Input(fmt.Sprintf("%s [Yn]", msg))

	return response == "Y" || response == "y"
}

// Outputf writes a formatted string to stdout
func (c *Context) Outputf(format string, args ...interface{}) {
	c.Output(fmt.Sprintf(format, args...))
}
