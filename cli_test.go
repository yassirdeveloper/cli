package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yassirdeveloper/cli/command"
	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
)

type mockOperator struct {
	output bytes.Buffer
}

func (m *mockOperator) Write(s string) errors.Error {
	_, err := m.output.Write([]byte(s))
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (m *mockOperator) String() string {
	return m.output.String()
}

func (m *mockOperator) Read() (string, errors.Error) {
	return "", nil
}

func TestNewCli(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	assert.Equal(t, "test-cli", cli.Name, "Name should match the provided name")
	assert.Equal(t, DEFAULT_SYMBOL, cli.Symbol, "Symbol should default to '>'")
	assert.Equal(t, DEFAULT_HISTORY_LIMIT, cli.HistoryLimit, "History limit should default to 100")
	assert.NotNil(t, cli.commander, "Commander should be initialized")
}

func TestInvalidNewCli(t *testing.T) {
	_, err := NewCli("test-cli", "bad")
	assert.Error(t, err, "an error should occur for invalid cli")
}

func TestSetVersion_Valid(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	version := "1.2.3"
	_, err = cli.SetVersion(version)

	assert.NoError(t, err, "No error should occur for valid version format")
	assert.Equal(t, "v"+version, command.GetVersionString(), "Version should be set correctly")
}

func TestSetVersion_Invalid(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	invalidVersion := "v1.2.t"
	_, err = cli.SetVersion(invalidVersion)

	assert.Error(t, err, "Error should occur for invalid version format")
	assert.Contains(t, err.Error(), "invalid version format", "Error message should indicate invalid format")
}

func TestAddCommand_Valid(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	cmd := command.NewCommand(
		"test",
		"A test command",
		func(input command.CommandInput, writer operator.Operator) errors.Error {
			err := writer.Write("Test command executed\n")
			return errors.NewUnexpectedError(err)
		},
	)

	err = cli.AddCommand(cmd)

	assert.NoError(t, err, "No error should occur for valid command")
	assert.Contains(t, cli.commander.GetCommands(), "test", "Command should be added to the commander")
}

func TestAddCommand_Invalid(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	cmd := command.NewCommand(
		"",
		"A test command",
		func(input command.CommandInput, writer operator.Operator) errors.Error {
			return nil
		},
	)

	err = cli.AddCommand(cmd)

	assert.Error(t, err, "Error should occur for invalid command")
	assert.Contains(t, err.Error(), "command name cannot be empty", "Error message should indicate invalid command")
}

func TestRun_NonInteractiveMode(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	// Add a mock command
	cli.AddCommand(
		command.NewCommand(
			"greet",
			"Greets the user",
			func(input command.CommandInput, writer operator.Operator) errors.Error {
				err := writer.Write("Hello, world!\n")
				return errors.NewUnexpectedError(err)
			},
		),
	)

	// Mock os.Args
	os.Args = []string{"cli", "greet"}

	var buf mockOperator
	cli.SetOperator(&buf)

	cli.Run(false)

	output := buf.String()
	assert.Contains(t, output, "Hello, world!", "Output should contain the greeting message")
}

func TestRun_InteractiveMode(t *testing.T) {
	cli, err := NewCli("test-cli", "0.0.0")
	assert.NoError(t, err, "No error should occur for valid cli")

	// Add a mock command
	cli.AddCommand(
		command.NewCommand(
			"echo",
			"Echoes the input",
			func(input command.CommandInput, writer operator.Operator) errors.Error {
				err := writer.Write(input.String() + "\n")
				return errors.NewUnexpectedError(err)
			},
		),
	)
	// Set up a custom writer to capture output
	var buf mockOperator
	cli.SetOperator(&buf)

	// Simulate interactive mode
	go func() {
		cli.Run(true)
	}()

	// TODO: Simulate user input (requires mocking readline or using a library like `os/exec`).
	// For now, this test ensures the Run method doesn't panic.
}
