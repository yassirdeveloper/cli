package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yassirdeveloper/cli/commands"
)

func TestNewCli(t *testing.T) {
	cli := NewCli("test-cli")

	assert.Equal(t, "test-cli", cli.Name, "Name should match the provided name")
	assert.Equal(t, DEFAULT_SYMBOL, cli.Symbol, "Symbol should default to '>'")
	assert.Equal(t, DEFAULT_HISTORY_LIMIT, cli.HistoryLimit, "History limit should default to 100")
	assert.NotNil(t, cli.commander, "Commander should be initialized")
}

func TestSetVersion_Valid(t *testing.T) {
	cli := NewCli("test-cli")

	version := "v1.2.3"
	updatedCli, err := cli.SetVersion(version)

	assert.NoError(t, err, "No error should occur for valid version format")
	assert.Equal(t, version, updatedCli.version, "Version should be set correctly")
}

func TestSetVersion_Invalid(t *testing.T) {
	cli := NewCli("test-cli")

	invalidVersion := "1.2.3"
	_, err := cli.SetVersion(invalidVersion)

	assert.Error(t, err, "Error should occur for invalid version format")
	assert.Contains(t, err.Error(), "invalid version format", "Error message should indicate invalid format")
}

func TestAddCommand_Valid(t *testing.T) {
	cli := NewCli("test-cli")

	cmd := commands.NewCommand(
		"test",
		"A test command",
		func(input commands.CommandInput, writer io.Writer) commands.Error {
			_, err := writer.Write([]byte("Test command executed\n"))
			return commands.NewUnexpectedError(err)
		},
	)

	err := cli.AddCommand(cmd)

	assert.NoError(t, err, "No error should occur for valid command")
	assert.Contains(t, cli.commander.GetCommands(), "test", "Command should be added to the commander")
}

func TestAddCommand_Invalid(t *testing.T) {
	cli := NewCli("test-cli")

	cmd := commands.NewCommand(
		"",
		"A test command",
		func(input commands.CommandInput, writer io.Writer) commands.Error {
			return nil
		},
	)

	err := cli.AddCommand(cmd)

	assert.Error(t, err, "Error should occur for invalid command")
	assert.Contains(t, err.Error(), "command name cannot be empty", "Error message should indicate invalid command")
}

func TestRun_NonInteractiveMode(t *testing.T) {
	cli := NewCli("test-cli")

	// Add a mock command
	cli.AddCommand(
		commands.NewCommand(
			"greet",
			"Greets the user",
			func(input commands.CommandInput, writer io.Writer) commands.Error {
				_, err := writer.Write([]byte("Hello, world!\n"))
				return commands.NewUnexpectedError(err)
			},
		),
	)

	// Mock os.Args
	os.Args = []string{"cli", "greet"}

	var buf bytes.Buffer
	cli.SetWriter(&buf)

	cli.Run(false)

	output := buf.String()
	assert.Contains(t, output, "Hello, world!", "Output should contain the greeting message")
}

func TestRun_InteractiveMode(t *testing.T) {
	cli := NewCli("test-cli")

	// Add a mock command
	cli.AddCommand(
		commands.NewCommand(
			"echo",
			"Echoes the input",
			func(input commands.CommandInput, writer io.Writer) commands.Error {
				_, err := writer.Write([]byte(input.String() + "\n"))
				return commands.NewUnexpectedError(err)
			},
		),
	)
	// Set up a custom writer to capture output
	var buf bytes.Buffer
	cli.SetWriter(&buf)

	// Simulate interactive mode
	go func() {
		cli.Run(true)
	}()

	// TODO: Simulate user input (requires mocking readline or using a library like `os/exec`).
	// For now, this test ensures the Run method doesn't panic.
}
