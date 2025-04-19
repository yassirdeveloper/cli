package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	// Set up a mock writer to capture output
	writer := &mockOperator{}

	// Create the VersionCommand
	version = "3.6.8"
	versionCommand := createVersionCommand(version)

	// Create a mock CommandInput (not used in this case, but required by the interface)
	input := &commandInput{
		arguments: make(map[string]any),
		options:   make(map[string]any),
	}

	t.Run("Valid Version Output", func(t *testing.T) {
		// Run the VersionCommand handler
		err := versionCommand.Handle(input, writer)
		assert.NoError(t, err, "Expected no error during execution")

		// Verify the output matches the expected version string
		expectedOutput := "v" + version
		assert.Equal(t, expectedOutput, writer.output.String(), "Output does not match expected version string")
	})
}
