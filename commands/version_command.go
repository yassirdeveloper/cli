package commands

import (
	"io"
)

var version string

func GetVersionString() string {
	return "v" + version
}

func versionHandler(_ CommandInput, writer io.Writer) Error {
	_, err := writer.Write([]byte(GetVersionString()))
	if err != nil {
		return &UnexpectedError{err: err}
	}
	return nil
}

func VersionCommand(v string) Command {
	command := &command{
		Name:        "version",
		Description: "Display the current version.",
	}
	version = v
	command.setHandler(versionHandler)
	return command
}
