package commands

import (
	"fmt"
	"io"
	"strings"
)

var helpText string

var commandOpt = commandOption{label: "command", letter: 'c', name: "command", valueType: TypeString, description: "Name of the command to get detailed help for"}

func helpHandler(input CommandInput, writer io.Writer) Error {
	commander := GetCommander()

	// If a specific command name is provided, show help for that command
	opt, err := input.ParseOption(commandOpt)
	if err != nil {
		return err
	}
	if opt != nil {
		cmdName := opt.(string)
		cmd, exists := commander.Get(cmdName)
		if exists {
			_, err := writer.Write([]byte("Command description:\n" + cmd.Help()))
			if err != nil {
				return &UnexpectedError{err: err}
			}
			return nil
		}
		_, err := writer.Write([]byte(fmt.Sprintf("No help available for command: %s\n", cmdName)))
		if err != nil {
			return &UnexpectedError{err: err}
		}
		return nil
	}

	// Otherwise, list help for all commands
	cmds := commander.GetCommands()
	var description strings.Builder
	if helpText != "" {
		description.WriteString(helpText)
	}
	description.WriteString("List of commands:\n")
	for _, cmd := range cmds {
		comm, exists := commander.Get(cmd)
		if !exists {
			panic(fmt.Sprintf("Commander does not return a Command for an existing command name %s", cmd))
		}
		description.WriteString(comm.Help())
	}
	_, err_ := writer.Write([]byte(description.String()))
	if err_ != nil {
		return &UnexpectedError{err: err_}
	}
	return nil
}

func HelpCommand(s string) Command {
	command := &command{
		Name:        "help",
		Description: "Display help information for commands.",
	}
	helpText = s
	command.AddOption(commandOpt)
	command.setHandler(helpHandler)
	return command
}
