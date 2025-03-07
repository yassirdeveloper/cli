package commands

import (
	"io"
	"os"
)

func exitHandler(_ CommandInput, _ io.Writer) Error {
	os.Exit(0)
	return nil
}

func ExitCommand() Command {
	command := &command{handler: exitHandler, helpText: "Exit the application."}
	return command
}
