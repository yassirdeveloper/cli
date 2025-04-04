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
	command := &command{Name: "exit", handler: exitHandler, Description: "Exit the application."}
	return command
}
