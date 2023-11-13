package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lkendrickd/cmdcraft/cmdcraft"
)

func main() {
	// Create a new CommandCraft instance
	commandCraft := cmdcraft.NewCommandCraft()

	// We are handling two example commands: echo and date (see below)
	// only the echo command has flags associated with it so we define them now

	// ### Command Number 1 our echo command #############################################################

	// Define flags for the echo command
	echoFlags := []cmdcraft.Flag{
		{LongOption: "message", ShortOption: "m", Type: "string", Description: "The message to echo"},
	}

	// Add the echo command to CommandCraft with the flags we defined above and the echo handler
	commandCraft.AddCommand(cmdcraft.NewCommand(
		"echo",
		"Echoes back the input",
		"echo --message <input>",
		echoFlags,            // The flags we defined above
		[]cmdcraft.Command{}, // No subcommands for the echo command
		echo,                 // The echo handler/function below which is our business logic function
	))

	// ### Command Number 2 our date command #############################################################

	// Define the date command without flags
	commandCraft.AddCommand(cmdcraft.NewCommand(
		"date",
		"Displays the current date",
		"date",
		[]cmdcraft.Flag{},    // No flags for the date command
		[]cmdcraft.Command{}, // No subcommands for the date command
		today,                // The today handler/function below which is our business logic function
	))

	// ### Execute CommandCraft ############################################################################
	// Execute CommandCraft with the arguments passed to the program
	// Examples:
	// go run cmd/main.go echo --help
	// go run cmd/main.go echo --message "Hello World"
	// go run cmd/main.go date
	// go run cmd/main.go --help

	err := commandCraft.Execute(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

/*
##################################################################
# Core Business Logic Functions - what you will be implementing  #
##################################################################
*/

// echo prints the values of the flags passed to it
func echo(data interface{}) error {
	cmd, ok := data.(*cmdcraft.Command)
	if !ok {
		return fmt.Errorf("invalid data type")
	}

	for key, value := range cmd.FlagValues {
		fmt.Printf("%s: %v\n", key, *value.(*string)) // Type assertion based on the flag type
	}
	return nil
}

// today prints the current date
func today(data interface{}) error {
	fmt.Printf("Today's date is %s\n", time.Now().Format("01-02-2006"))
	return nil
}
