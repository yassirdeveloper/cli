package command

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	exitCmd := createExitCommand()
	versionCmd := createVersionCommand("0.0.0")
	commander := GetCommander()
	commander.AddCommand("exit", exitCmd)
	commander.AddCommand("version", versionCmd)

	helpCommand := createHelpCommand()

	writer := &mockOperator{}

	t.Run("List All Commands", func(t *testing.T) {
		input := &commandInput{
			arguments: map[string]any{},
			options:   map[string]any{},
		}

		err := helpCommand.Handle(input, writer)
		assert.NoError(t, err)
		assert.NotEmpty(t, writer.String())
	})

	t.Run("Get Help for Specific Command", func(t *testing.T) {
		writer.Reset()

		input := &commandInput{
			arguments: map[string]any{},
			options: map[string]any{
				"command": "exit",
			},
		}

		err := helpCommand.Handle(input, writer)
		assert.NoError(t, err)
		assert.NotEmpty(t, strings.TrimSpace(writer.String()))
	})

	t.Run("Nonexistent Command", func(t *testing.T) {
		writer.Reset()

		input := &commandInput{
			arguments: map[string]any{},
			options: map[string]any{
				"command": "nonexistent",
			},
		}

		err := helpCommand.Handle(input, writer)
		assert.NoError(t, err)

		expectedOutput := "No help available for command: nonexistent\n"
		assert.Equal(t, expectedOutput, writer.String())
	})
}
