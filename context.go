package reaper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Context is used to hold the state of the application
type Context struct {
	app       *App
	args      []string
	flags     map[string]interface{}
	namedArgs map[string]string
}

func newContext(app *App, args []string) *Context {
	return &Context{
		app:       app,
		args:      args,
		flags:     make(map[string]interface{}, 0),
		namedArgs: make(map[string]string, 0),
	}
}

// Flag returns a flag with the passed name. If the flag is not found an error will be returned.
func (c *Context) Flag(name string) (interface{}, error) {
	if val, ok := c.flags[name]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("missing required flag %s", name)
}

// FlagString returns a flag string value.
//
// This will panic if the flag does not exist
func (c *Context) FlagString(name string) string {
	value, _ := c.Flag(name)
	return *(value.(*string))
}

// Argument returns a named argument
func (c *Context) Argument(name string) string {
	return c.namedArgs[name]
}

// Output writes the string to stdout
func (c *Context) Output(value string) {
	c.app.output.Print(value)
}

// Get asks for user input
func (c *Context) Get(msg string) (string, error) {
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
	response, _ := c.Get(fmt.Sprintf("%s [Yn]", msg))

	return response == "Y" || response == "y"
}

// Outputf writes a formatted string to stdout
func (c *Context) Outputf(format string, args ...interface{}) {
	c.Output(fmt.Sprintf(format, args...))
}
