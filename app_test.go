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

	t.Run("given a valid app with an unknown command", func(t *testing.T) {
		app := NewApp("testing")

		app.Command("greet", func(c *Context) error {
			return nil
		})

		err := app.Execute([]string{"foo"})

		fmt.Println(err)

		assert.Error(t, err)
	})
}
