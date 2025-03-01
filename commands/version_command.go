package commands

import (
	"io"
)

const version = "0.0.0"

func versionHandler(_ CommandInput, writer io.Writer) Error {
	_, err := writer.Write([]byte("v" + version))
	if err != nil {
		return &UnexpectedError{}
	}
	return nil
}

func VersionCommand() Command {
	command := &command{handler: versionHandler}
	return command
}
