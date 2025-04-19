package operator

import (
	"bufio"
	"io"
	"os"

	"github.com/yassirdeveloper/cli/errors"
)

type Operator interface {
	Write(string) errors.Error
	Read() (string, errors.Error)
}

type stdOperator struct {
	delim       byte
	maxReadSize int
	writer      io.Writer
	reader      *bufio.Reader
}

func (o *stdOperator) Write(s string) errors.Error {
	_, err := o.writer.Write([]byte(s))
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (o *stdOperator) Read() (string, errors.Error) {
	s, err := o.reader.ReadString(o.delim)
	if err != nil {
		return s, errors.NewUnexpectedError(err)
	}
	return s, nil
}

func NewStdOperator(delim byte, maxReadSize int) *stdOperator {
	return &stdOperator{
		delim:       delim,
		maxReadSize: maxReadSize,
		writer:      os.Stdout,
		reader:      bufio.NewReader(os.Stdin),
	}
}
