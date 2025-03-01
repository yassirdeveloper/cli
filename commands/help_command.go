package commands

import (
	"io"
)

func helpHandler(input CommandInput, writer io.Writer) Error {
	_, err := writer.Write([]byte("help"))
	if err != nil {
		return &UnexpectedError{}
	}
	return nil
}

func HelpCommand() Command {
	command := &command{handler: helpHandler}
	return command
}
