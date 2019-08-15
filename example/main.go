package main

import (
	"errors"
	"log"
	"os"

	"github.com/maddiesch/reaper"
)

func main() {
	app := reaper.NewApp("example")

	app.Description = "This is an example app for the Reaper CLI library"

	app.Command("greet", func(c *reaper.Context) error {
		name, err := c.Input("What's your name?")
		if err != nil {
			return err
		}

		c.Outputf("%s %s!\n", c.FlagString("greeting"), name)

		return nil
	}).Configure(func(c *reaper.Command) {
		c.Description = "prints the greeting and name"
		c.Flag("greeting", "string", "Hi", "the greeting to perform")
	})

	app.Command("get", func(c *reaper.Context) error {
		switch c.Argument("key") {
		case "private/information":
			if c.Argument("secret") != "sup3rsekret" {
				return errors.New("Invalid Secret")
			}
			if c.Confirm("are you sure?") {
				c.Output("Secret Private Data")
			}
		case "public/information":
			c.Output("Public Data")
		default:
			return errors.New("Unknown key")
		}

		return nil
	}).Configure(func(c *reaper.Command) {
		c.Description = "fetch a value from the thing"
		c.Argument("key", true, "the key for the thing to fetch")
		c.Argument("secret", false, "if the key is protected, the secret used to unlock it")
		c.Flag("port", "integer", int64(80), "the port of the thing")

		c.Example("-port 442 private/information sup3rsekret")
		c.Example("public/information")
	})

	err := app.Execute(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
