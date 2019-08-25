package main

import (
	"os"

	"github.com/maddiesch/reaper"
)

func main() {
	app := reaper.NewApp("reaper")
	app.Version = reaper.Version
	app.Description = "A collection of commands that makes working with reaper easier"

	app.Command("new", newAppHandler).Configure(func(c *reaper.Command) {
		c.Description = "creates a new reaper application from a template"

		c.Flag("type", "string", "cli", "the type of cli app to create")
		c.Flag("github", "string", "ghusername", "the name of the user for module support")
		c.Flag("path", "string", ".", "path for the new application")

		c.Argument("name", true, "the name of the application")

		c.Example("cli-app")
	})

	err := app.Execute(os.Args[1:])
	if err != nil {
		reaper.Fatal(err)
	}
}
