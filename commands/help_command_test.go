package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a sample HelpCommand with a mock commander
func createHelpCommand() Command {
	helpCmd := HelpCommand()
	return helpCmd
}

func TestHelpCommand(t *testing.T) {
	exitCmd := &command{
		helpText: "Exit the application.",
	}
	versionCmd := &command{
		helpText: "Display the application version.",
	}
	commander := GetCommander()
	commander.AddCommand("exit", exitCmd)
	commander.AddCommand("version", versionCmd)

	helpCommand := createHelpCommand()

	// Mock writer to capture output
	writer := &mockWriter{}

	t.Run("List All Commands", func(t *testing.T) {
		input := &commandInput{
			arguments: map[string]any{},
			options:   map[string]any{},
		}

		err := helpCommand.Handle(input, writer)
		assert.NoError(t, err)

		expectedOutput := `-exit:Exit the application.
Usage: exit
-version:Display the application version.
Usage: version
`
		assert.Equal(t, expectedOutput, writer.String())
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

		expectedOutput := "Exit the application.\nUsage: exit"
		assert.Equal(t, expectedOutput, strings.TrimSpace(writer.String()))
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
