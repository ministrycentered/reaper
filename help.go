package reaper

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

func helpCommandHandler(c *Context) error {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	output := log.New(writer, "", 0)

	output.Print(aurora.Bold(fmt.Sprintf("%s", c.app.fullName())))
	output.Print(aurora.Faint(fmt.Sprintf("Version: %s", c.app.Version)))
	output.Print("")

	if c.app.Description != "" {
		output.Print(c.app.Description)
		output.Print("")
	}

	output.Print(aurora.Bold(aurora.Magenta("Commands")))
	for _, name := range c.app.commandNames() {
		cmd := c.app.commands[name]
		var details strings.Builder

		details.WriteString("")
		details.WriteString(aurora.White(name).String())
		if cmd.Description != "" {
			details.WriteString(" -- ")
			details.WriteString(cmd.Description)
		}

		output.Print(details.String())

		if len(cmd.args) > 0 {
			output.Print("  ", aurora.Faint("Arguments:"))
			for _, arg := range cmd.args {
				var fName string
				if arg.required {
					fName = aurora.Cyan(arg.name).String()
				} else {
					fName = aurora.Cyan(fmt.Sprintf("[%s]", arg.name)).String()
				}
				output.Printf("    %s -- %s", fName, arg.usage)
			}
		}

		if len(cmd.flags) > 0 {
			output.Print("  ", aurora.Faint("Flags:"))
			for _, name := range cmd.flagNames() {
				info := cmd.flags[name]
				fName := aurora.Cyan(fmt.Sprintf("-%s", name))
				fType := aurora.Faint(fmt.Sprintf("[%s]", info.kind))
				output.Printf("    %s %s -- %s (default `%v`)", fName, fType, info.usage, info.value)
			}
		}

		if len(cmd.examples) > 0 {
			output.Print("  ", aurora.Faint("Examples:"))
			for _, example := range cmd.examples {
				output.Printf("      $ %s %s %s", c.app.name, name, example)
			}
		}

		output.Print("")
	}

	appNames := c.app.appNames()
	if len(appNames) > 0 {
		output.Print(aurora.Yellow("Apps:"), aurora.Faint("Sub applications for"), aurora.Faint(c.app.name))

		for _, name := range appNames {
			output.Print(" - ", aurora.Cyan(name))
		}

		output.Print("")
	}

	err := writer.Flush()
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stderr, buffer.String())

	return err
}
