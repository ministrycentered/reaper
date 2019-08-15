package reaper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	t.Run("given valid inputs for a valid application", func(t *testing.T) {
		app := NewApp("testing")

		var name string

		app.Command("greet", func(c *Context) error {
			name = c.FlagString("name")
			return nil
		}).Configure(func(c *Command) {
			c.Flag("name", "string", "human", "the name of the person to greet")
		})

		err := app.Execute([]string{"greet", "-name", "Maddie"})

		assert.NoError(t, err)
		assert.Equal(t, "Maddie", name)
	})

	t.Run("given a valid app with a before action", func(t *testing.T) {
		app := NewApp("testing")

		var value string

		app.Before(func(c *Context) error {
			c.Set("value", "Super Cool")
			return nil
		})

		app.Command("greet", func(c *Context) error {
			value = c.Get("value").(string)
			return nil
		})

		err := app.Execute([]string{"greet"})

		assert.NoError(t, err)
		assert.Equal(t, "Super Cool", value)
	})

	t.Run("given a valid app with an unknown command", func(t *testing.T) {
		app := NewApp("testing")

		app.Before(func(c *Context) error {
			c.Set("foo", "bar")
			return nil
		})

		app.Command("greet", func(c *Context) error {
			assert.Equal(t, "bar", c.Get("foo"))
			return nil
		})

		err := app.Execute([]string{"foo"})

		fmt.Println(err)

		assert.Error(t, err)
	})
}
