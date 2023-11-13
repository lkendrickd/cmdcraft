package cmdcraft

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

// CommandCraft is the primary controller for the CLI
type CommandCraft struct {
	Commands []Command
}

// NewCommandCraft creates a new CommandCraft instance
func NewCommandCraft() *CommandCraft {
	return &CommandCraft{
		Commands: []Command{},
	}
}

// CommandCraftHelp prints the help menu for CommandCraft
func (c *CommandCraft) CommandCraftHelp() error {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("#============================#\n")
	sb.WriteString("#   CommandCraft Help Menu   #\n")
	sb.WriteString("#============================#\n\n")

	for _, cmd := range c.Commands {
		sb.WriteString(fmt.Sprintf("Command: %s\n", cmd.Name))
		sb.WriteString(fmt.Sprintf("Description: %s\n", cmd.Description))
		sb.WriteString(fmt.Sprintf("Usage: %s\n", cmd.Usage))

		if len(cmd.Subcommands) > 0 {
			sb.WriteString("Subcommands:\n")
			for _, subCmd := range cmd.Subcommands {
				sb.WriteString(fmt.Sprintf("  - %s: %s\n", subCmd.Name, subCmd.Description))
			}
		}
		sb.WriteString("\n")
	}

	log.Println(sb.String())
	return nil
}

// AddCommand adds a command to CommandCraft
func (c *CommandCraft) AddCommand(cmd Command) {
	// add the new command
	c.Commands = append(c.Commands, cmd)
}

// Execute executes the command
func (c *CommandCraft) Execute(args []string) error {
	if len(args) <= 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return c.CommandCraftHelp()
	}

	// Find the main command
	var mainCmd *Command
	for _, cmd := range c.Commands {
		if cmd.Name == args[0] {
			mainCmd = &cmd
			break
		}
	}
	if mainCmd == nil {
		return fmt.Errorf("command not found")
	}

	// Check if there's a subcommand
	if len(args) > 1 {
		if args[1] == "help" || args[1] == "--help" || args[1] == "-h" {
			mainCmd.CommandHelp()
			return nil
		}

		for _, subCmd := range mainCmd.Subcommands {
			if subCmd.Name == args[1] {
				// Found a subcommand, now handle it
				return handleSubCommand(&subCmd, args[2:])
			}
		}
	}

	// If no subcommand, process the main command
	flagSet := flag.NewFlagSet(mainCmd.Name, flag.ContinueOnError)
	mainCmd.initFlags(flagSet)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}
	return mainCmd.Handler(mainCmd)
}

// handleSubCommand handles the execution of a subcommand
func handleSubCommand(subCmd *Command, args []string) error {
	flagSet := flag.NewFlagSet(subCmd.Name, flag.ContinueOnError)
	subCmd.initFlags(flagSet)
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	return subCmd.Handler(subCmd)
}
