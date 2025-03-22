package commands

import (
	"fmt"
	"io"
	"strings"
)

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
			_, err := writer.Write([]byte(cmd.Help()))
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
	var helpText strings.Builder
	helpText.WriteString("List of commands:\n")
	for _, cmd := range cmds {
		comm, exists := commander.Get(cmd)
		if !exists {
			panic(fmt.Sprintf("Commander does not return a Command for an existing command name %s", cmd))
		}
		helpText.WriteString(fmt.Sprintf("\t- %-15s %s\n", cmd+":", comm.Help()))
	}
	_, err_ := writer.Write([]byte(helpText.String()))
	if err_ != nil {
		return &UnexpectedError{err: err_}
	}
	return nil
}

func HelpCommand() Command {
	command := &command{
		helpText: "Display help information for commands.",
	}
	command.AddOption(commandOpt)
	command.setHandler(helpHandler)
	return command
}
