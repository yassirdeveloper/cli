package command

import (
	"fmt"
	"strings"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
)

var helpText string

var commandOpt = CommandOption{
	Label:  "command",
	Letter: 'c', Name: "command",
	ValueType:   TypeString,
	Description: "Name of the command to get detailed help for",
}

func helpHandler(input CommandInput, operator operator.Operator) errors.Error {
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
			err := operator.Write("Command description:\n" + cmd.Help())
			if err != nil {
				return errors.NewUnexpectedError(err)
			}
			return nil
		}
		err := operator.Write(fmt.Sprintf("No help available for command: %s\n", cmdName))
		if err != nil {
			return errors.NewUnexpectedError(err)
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
	err_ := operator.Write(description.String())
	if err_ != nil {
		return errors.NewUnexpectedError(err_)
	}
	return nil
}

func HelpCommand(s string) Command {
	cmd := NewCommand(
		"help",
		"Display help information for commands.",
		helpHandler,
	)
	helpText = s
	cmd.AddOption(commandOpt)
	return cmd
}
