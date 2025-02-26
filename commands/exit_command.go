package commands

import (
	"io"
	"os"
)

func exitHandler(_ CommandInput, _ io.Writer) Error {
	os.Exit(0)
	return nil
}

var ExitCommand = NewCommand("exit").SetHandler(exitHandler)
