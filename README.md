# Reaper

A go library for creating a command line app.

```go
func main() {
  app := reaper.NewApp("example")

  app.Description = "This is an example app for the Reaper CLI library"

  app.Command("greet", func(c *reaper.Context) error {
    name, err := c.Get("What's your name?")
    if err != nil {
      return err
    }

    c.Outputf("%s %s!\n", c.FlagString("greeting"), name)

    return nil
  }).Configure(func(c *reaper.Command) {
    c.Description = "prints the greeting and name"
    c.Flag("greeting", "string", "Hi", "the greeting to perform")
  })

  err := app.Execute(os.Args[1:])
  if err != nil {
    log.Fatal(err)
  }
}
```

`$ example greet`

## Builtin Help

Given the above example, out of the box the `help` sub-command will print this.

```text
example help
Version: 1.0.0

This is an example app for the Reaper CLI library

Commands
greet -- prints the greeting and name
  Flags:
    -greeting [string] -- the greeting to perform (default Hi)

help -- prints this help dialogue

version -- print the current version
```
