package command

import (
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
)

var version string

func GetVersionString() string {
	return "v" + version
}

func versionHandler(_ CommandInput, operator operator.Operator) errors.Error {
	err := operator.Write(GetVersionString())
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func VersionCommand(v string) Command {
	cmd := NewCommand(
		"version",
		"Display the current version.",
		versionHandler,
	)
	version = v
	return cmd
}
