package commands

import (
	"fmt"
)

type Error interface {
	Error() string
	Display() string
}

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
	command string
}

func (e *InvalidCommandUsageError) Error() string {
	return fmt.Sprintf("Invalid usage of command: %s", e.command)
}

func (e *InvalidCommandUsageError) Display() string {
	return fmt.Sprintf("Invalid usage of command: %s", e.command)
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

func (e *CommandError) Error() string {
	return e.message
}

func (e *CommandError) Display() string {
	return e.message
}

type UnexpectedError struct {
	message string
	err     error
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}

func (e *UnexpectedError) Display() string {
	return "An unexpected error occured!"
}

type SetupError struct {
	message string
}

func (e *SetupError) Error() string {
	return e.message
}

func (e *SetupError) Display() string {
	return "An error occured during steup:" + e.message
}
