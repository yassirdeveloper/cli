package commands

import (
	"io"
)

const version = "0.0.0"

func versionHandler(input CommandInput, writer io.Writer) Error {
	_, err := writer.Write([]byte("v" + version))
	if err != nil {
		return &UnexpectedError{}
	}
	return nil
}

var VersionCommand = NewCommand("version").SetHandler(versionHandler)
