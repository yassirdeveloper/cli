package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
)

func createExitCommand() Command {
	exitCmd := ExitCommand()
	return exitCmd
}

func createVersionCommand(v string) Command {
	return VersionCommand(v)
}

func createHelpCommand() Command {
	helpCmd := HelpCommand("")
	return helpCmd
}

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

func (m *mockOperator) Reset() {
	m.output.Reset()
}

func (m *mockOperator) String() string {
	return m.output.String()
}

func (m *mockOperator) Read() (string, errors.Error) {
	return "", nil
}

// Helper function to create a sample command
func createSampleCommand() Command {
	comm := &command{}
	comm.setName("test")
	var err errors.Error

	// Add argument
	if _, err = comm.AddArgument(CommandArgument{
		Label:       "arg1",
		Description: "First argument",
		Position:    0,
		ValueType:   TypeString,
	}); err != nil {
		panic(err)
	}

	// Add option
	if _, err = comm.AddOption(CommandOption{
		Label:       "opt1",
		Description: "First option",
		Letter:      'o',
		Name:        "option1",
		ValueType:   TypeString,
	}); err != nil {
		panic(err)
	}

	// Set handler
	comm.setHandler(func(input CommandInput, operator operator.Operator) errors.Error {
		argValue, err := input.ParseArgument(CommandArgument{Label: "arg1", ValueType: TypeString})
		if err != nil {
			return err
		}
		optValue, _ := input.ParseOption(CommandOption{Label: "opt1", ValueType: TypeString})
		operator.Write("Arg: " + argValue.(string) + ", Opt: " + optValue.(string))
		return nil
	})

	return comm
}

func TestCommandParsing(t *testing.T) {
	comm := createSampleCommand()

	t.Run("Valid Input", func(t *testing.T) {
		input, err := comm.Parse([]string{"value1", "-o", "optionValue"})
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		if input == nil {
			t.Fatal("Parsed input should not be nil")
		}

		argValue, err := input.ParseArgument(CommandArgument{Label: "arg1", ValueType: TypeString})
		if err != nil {
			t.Fatalf("Expected no error parsing argument, but got: %v", err)
		}
		if argValue != "value1" {
			t.Errorf("Expected argument value 'value1', but got '%v'", argValue)
		}

		optValue, err := input.ParseOption(CommandOption{Label: "opt1", ValueType: TypeString})
		if err != nil {
			t.Fatalf("Expected no error parsing option, but got: %v", err)
		}
		if optValue != "optionValue" {
			t.Errorf("Expected option value 'optionValue', but got '%v'", optValue)
		}
	})

	t.Run("Missing Argument", func(t *testing.T) {
		_, err := comm.Parse([]string{})
		if err == nil {
			t.Fatal("Expected an error for missing argument, but got none")
		}
		if _, ok := err.(*InvalidCommandUsageError); !ok {
			t.Errorf("Expected InvalidCommandUsageError, but got %T", err)
		}
	})

	t.Run("Invalid Option", func(t *testing.T) {
		_, err := comm.Parse([]string{"value1", "-x", "invalidOption"})
		if err == nil {
			t.Fatal("Expected an error for invalid option, but got none")
		}
		if _, ok := err.(*UnreconizedFlagError); !ok {
			t.Errorf("Expected UnreconizedFlagError, but got %T", err)
		}
	})
}

func TestCommandExecution(t *testing.T) {
	comm := createSampleCommand()
	writer := &mockOperator{}

	t.Run("Successful Execution", func(t *testing.T) {
		err := comm.Handle(&commandInput{
			arguments: map[string]any{
				"arg1": "testArg",
			},
			options: map[string]any{
				"opt1": "testOpt",
			},
		}, writer)
		if err != nil {
			t.Fatalf("Expected no error during execution, but got: %v", err)
		}
		if strings.TrimSpace(writer.output.String()) != "Arg: testArg, Opt: testOpt" {
			t.Errorf("Unexpected output: %s", writer.output.String())
		}
	})

	t.Run("Handler Error", func(t *testing.T) {
		failingCommand := &command{}
		failingCommand.setHandler(func(input CommandInput, writer operator.Operator) errors.Error {
			return &CommandError{message: "handler failed"}
		})

		err := failingCommand.Handle(nil, writer)
		if err == nil || err.Error() != "handler failed" {
			t.Errorf("Expected error 'handler failed', but got: %v", err)
		}
	})
}

func TestCommander(t *testing.T) {
	commander := GetCommander()
	writer := &mockOperator{}

	t.Run("Add and Get Command", func(t *testing.T) {
		command := createSampleCommand()
		commander.AddCommand("test", command)

		retrievedCommand, exists := commander.Get("test")
		if !exists {
			t.Fatalf("Expected retrieving command, but got none")
		}
		if retrievedCommand.String() != "test" {
			t.Errorf("Expected command name 'test', but got '%s'", retrievedCommand.String())
		}
	})

	t.Run("Get Nonexistent Command", func(t *testing.T) {
		_, exists := commander.Get("nonexistent")
		if exists {
			t.Fatal("Expected to not get a command for nonexistent command, but got one")
		}
	})

	t.Run("Run Command", func(t *testing.T) {
		commander.AddCommand("runTest", createSampleCommand())
		commander.SetOperator(writer)

		err := commander.Run([]string{"runTest", "argValue", "-o", "optValue"})
		if err != nil {
			t.Fatalf("Expected no error running command, but got: %v", err)
		}
		if strings.TrimSpace(writer.output.String()) != "Arg: argValue, Opt: optValue" {
			t.Errorf("Unexpected output: %s", writer.output.String())
		}
	})

	t.Run("Run Invalid Command", func(t *testing.T) {
		err := commander.Run([]string{"invalidCmd"})
		if err == nil {
			t.Fatal("Expected an error for invalid command, but got none")
		}
		if _, ok := err.(*InvalidCommandError); !ok {
			t.Errorf("Expected InvalidCommandError, but got %T", err)
		}
	})

	t.Run("Run Command with Missing Arguments", func(t *testing.T) {
		err := commander.Run([]string{"runTest"})
		if err == nil {
			t.Fatal("Expected an error for missing arguments, but got none")
		}
		if _, ok := err.(*InvalidCommandUsageError); !ok {
			t.Errorf("Expected InvalidCommandUsageError, but got %T", err)
		}
	})
}

func TestDuplicateArgumentsAndOptions(t *testing.T) {
	comm := &command{}
	comm.setName("test")

	t.Run("Duplicate Argument", func(t *testing.T) {
		_, err := comm.AddArgument(CommandArgument{
			Label:       "arg1",
			Description: "First argument",
			Position:    0,
			ValueType:   TypeString,
		})
		if err != nil {
			t.Fatalf("Unexpected error adding first argument: %v", err)
		}

		_, err = comm.AddArgument(CommandArgument{
			Label:       "arg1",
			Description: "Duplicate argument",
			Position:    1,
			ValueType:   TypeString,
		})
		if err == nil {
			t.Fatal("Expected an error for duplicate argument, but got none")
		}
		if _, ok := err.(*errors.SetupError); !ok {
			t.Errorf("Expected SetupError, but got %T", err)
		}
	})

	t.Run("Duplicate Option", func(t *testing.T) {
		_, err := comm.AddOption(CommandOption{
			Label:       "opt1",
			Description: "First option",
			Letter:      'o',
			Name:        "option1",
			ValueType:   TypeString,
		})
		if err != nil {
			t.Fatalf("Unexpected error adding first option: %v", err)
		}

		_, err = comm.AddOption(CommandOption{
			Label:       "opt1",
			Description: "Duplicate option",
			Letter:      'o',
			Name:        "option1",
			ValueType:   TypeString,
		})
		if err == nil {
			t.Fatal("Expected an error for duplicate option, but got none")
		}
		if _, ok := err.(*errors.SetupError); !ok {
			t.Errorf("Expected SetupError, but got %T", err)
		}
	})
}
