package cmdcraft

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

// Command is a single command that can be executed by CommandCraft
type Command struct {
	Name        string
	Description string
	Usage       string
	Subcommands []Command
	Flags       []Flag
	InitFlags   func(*flag.FlagSet)
	Handler     func(data interface{}) error
	FlagValues  map[string]interface{}
}

// NewCommand creates a new Command
func NewCommand(
	name string,
	description string,
	usage string,
	flags []Flag,
	subcommands []Command,
	handler func(data interface{}) error) Command {
	return Command{
		Name:        name,
		Description: description,
		Usage:       usage,
		Subcommands: subcommands,
		Flags:       flags,
		Handler:     handler,
	}
}

// Flag represents a flag that can be passed to a command
type Flag struct {
	LongOption  string
	ShortOption string
	Type        string
	Description string
}

// CommandHelp prints the help menu for a Command
func (c *Command) CommandHelp() error {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Command: %s\n", c.Name))
	sb.WriteString(fmt.Sprintf("Description: %s\n", c.Description))
	sb.WriteString(fmt.Sprintf("Usage: %s\n", c.Usage))

	if len(c.Subcommands) > 0 {
		sb.WriteString("Subcommands:\n")
		for _, subCmd := range c.Subcommands {
			sb.WriteString(fmt.Sprintf("  - %s: %s\n", subCmd.Name, subCmd.Description))
		}
	}

	log.Println(sb.String())
	return nil
}

// initFlags initializes the flags for a Command
func (c *Command) initFlags(flagSet *flag.FlagSet) {
	c.FlagValues = make(map[string]interface{})

	for _, f := range c.Flags {
		switch f.Type {
		case "string":
			var s string
			flagSet.StringVar(&s, f.LongOption, "", f.Description)
			flagSet.StringVar(&s, f.ShortOption, "", f.Description+" (short)")
			c.FlagValues[f.LongOption] = &s
		case "int":
			var i int
			flagSet.IntVar(&i, f.LongOption, 0, f.Description)
			flagSet.IntVar(&i, f.ShortOption, 0, f.Description+" (short)")
			c.FlagValues[f.LongOption] = &i
		case "bool":
			var b bool
			flagSet.BoolVar(&b, f.LongOption, false, f.Description)
			flagSet.BoolVar(&b, f.ShortOption, false, f.Description+" (short)")
			c.FlagValues[f.LongOption] = &b

		case "float64":
			var f64 float64
			flagSet.Float64Var(&f64, f.LongOption, 0, f.Description)
			flagSet.Float64Var(&f64, f.ShortOption, 0, f.Description+" (short)")
			c.FlagValues[f.LongOption] = &f64

		case "duration":
			var d time.Duration
			flagSet.DurationVar(&d, f.LongOption, 0, f.Description)
			flagSet.DurationVar(&d, f.ShortOption, 0, f.Description+" (short)")
			c.FlagValues[f.LongOption] = &d

		default:
			panic(fmt.Sprintf("invalid flag type: %s", f.Type))
		}
	}
}
