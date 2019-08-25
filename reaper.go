package reaper

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// Fatal prints the error and exits with a code of 1
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", aurora.Red("Fatal:"), err)
	if re, ok := err.(Error); ok {
		os.Exit(re.ExitStatus)
	}
	os.Exit(1)
}
