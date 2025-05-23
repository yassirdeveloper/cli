package command

import (
	"fmt"
)

type InvalidCommandError struct {
	command string
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("Invalid command: %s", e.command)
}

func (e *InvalidCommandError) Display() string {
	return fmt.Sprintf("Invalid command: %s", e.command)
}

type InvalidCommandUsageError struct {
	command Command
}

func (e *InvalidCommandUsageError) Error() string {
	return fmt.Sprintf("Invalid usage of command: %s", e.command.String())
}

func (e *InvalidCommandUsageError) Display() string {
	commandName := e.command.String()
	return fmt.Sprintf("Invalid usage of command: %s\n\n> %s: %s\n", commandName, commandName, e.command.Help())
}

type UnreconizedFlagError struct {
	command string
	flag    string
}

func (e *UnreconizedFlagError) Error() string {
	return fmt.Sprintf("Unreconized flag %s for command %s", e.flag, e.command)
}

func (e *UnreconizedFlagError) Display() string {
	return fmt.Sprintf("Unreconized flag %s for command %s", e.flag, e.command)
}

type CommandError struct {
	message string
}

func NewUCommandError(message string) *CommandError {
	return &CommandError{message: message}
}

func (e *CommandError) Error() string {
	return e.message
}

func (e *CommandError) Display() string {
	return e.message
}
