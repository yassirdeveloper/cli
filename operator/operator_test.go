package operator

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yassirdeveloper/cli/errors"
)

// Test Write method
func TestWrite(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a custom stdOperator with the buffer as the writer
	op := &stdOperator{
		writer: &buf,
	}

	// Write a string
	err := op.Write("Hello, World!")
	assert.NoError(t, err, "Write should not return an error")

	// Verify the written content
	assert.Equal(t, "Hello, World!", buf.String(), "Written content should match")
}

// Test Read method with successful input
func TestRead_Success(t *testing.T) {
	// Create a reader with predefined input
	input := "Test Input\n"
	reader := bufio.NewReader(strings.NewReader(input))

	// Create a custom stdOperator with the reader
	op := &stdOperator{
		delim:  '\n',
		reader: reader,
	}

	// Read the input
	result, err := op.Read()
	assert.NoError(t, err, "Read should not return an error")
	assert.Equal(t, "Test Input\n", result, "Read content should match the input")
}

// Test Read method with EOF (end of file)
func TestRead_EOF(t *testing.T) {
	// Create a reader with no input
	reader := bufio.NewReader(strings.NewReader(""))

	// Create a custom stdOperator with the reader
	op := &stdOperator{
		delim:  '\n',
		reader: reader,
	}

	// Attempt to read
	result, err := op.Read()
	assert.Error(t, err, "Read should return an error for EOF")
	assert.True(t, errors.IsUnexpectedError(err), "Error should be an UnexpectedError")
	assert.Equal(t, "", result, "Result should be empty for EOF")
}

// Test Read method with an unexpected error
func TestRead_UnexpectedError(t *testing.T) {
	// Create a failing reader that always returns an error
	failingReader := &failingReader{err: errors.New("mock read error")}

	// Create a custom stdOperator with the failing reader
	op := &stdOperator{
		delim:  '\n',
		reader: bufio.NewReader(failingReader),
	}

	// Attempt to read
	result, err := op.Read()
	assert.Error(t, err, "Read should return an error for a failing reader")
	assert.True(t, errors.IsUnexpectedError(err), "Error should be an UnexpectedError")
	assert.Equal(t, "", result, "Result should be empty for a failing reader")
}

// Helper: A failing reader that always returns an error
type failingReader struct {
	err error
}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}
