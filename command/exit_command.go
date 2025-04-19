package command

import (
	"os"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
)

func exitHandler(_ CommandInput, _ operator.Operator) errors.Error {
	os.Exit(0)
	return nil
}

func ExitCommand() Command {
	return NewCommand(
		"exit",
		"Exit the application.",
		exitHandler,
	)
}
