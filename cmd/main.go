package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lkendrickd/cmdcraft/cmdcraft"
)

const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitUsageError   = 2
)

func main() {
	commandCraft := cmdcraft.NewCommandCraft()

	// Define flags for the echo command
	echoFlags := []cmdcraft.Flag{
		{LongOption: "message", ShortOption: "m", Type: "string", Description: "The message to echo"},
	}

	// Add the echo command
	commandCraft.AddCommand(cmdcraft.NewCommand(
		"echo",
		"Echoes back the input",
		"echo --message <input>",
		echoFlags,
		[]cmdcraft.Command{},
		echo,
	))

	// Add the date command
	commandCraft.AddCommand(cmdcraft.NewCommand(
		"date",
		"Displays the current date",
		"date",
		[]cmdcraft.Flag{},
		[]cmdcraft.Command{},
		today,
	))

	// Execute
	if err := commandCraft.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(ExitGeneralError)
	}
}

// echo prints the message flag value
func echo(data interface{}) error {
	cmd, ok := data.(*cmdcraft.Command)
	if !ok {
		return fmt.Errorf("invalid data type: expected *Command")
	}

	message := *cmd.FlagValues["message"].(*string)
	if message == "" {
		return fmt.Errorf("--message is required")
	}

	fmt.Println(message)
	return nil
}

// today prints the current date
func today(data interface{}) error {
	fmt.Printf("Today's date is %s\n", time.Now().Format("2006-01-02"))
	return nil
}
