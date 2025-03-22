package commands

import (
	"io"
)

var Version = "0.0.0"

func getVersionString() string {
	return "v" + Version
}

func versionHandler(_ CommandInput, writer io.Writer) Error {
	_, err := writer.Write([]byte(getVersionString()))
	if err != nil {
		return &UnexpectedError{err: err}
	}
	return nil
}

func VersionCommand() Command {
	command := &command{
		helpText: "Display the current version.",
	}
	command.setHandler(versionHandler)
	return command
}
