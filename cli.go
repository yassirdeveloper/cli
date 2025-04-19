package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	readline "github.com/chzyer/readline"
	"github.com/yassirdeveloper/cli/command"
	"github.com/yassirdeveloper/cli/operator"
)

const DEFAULT_SYMBOL = ">"
const DEFAULT_HISTORY_LIMIT = 100

var DEFAULT_OPERATOR = operator.NewStdOperator('\n', 4096)
var cliInstance *Cli

type Cli struct {
	Name         string
	HistoryLimit int
	Symbol       string
	commander    command.Commander
}

func GetCliInstance() *Cli {
	return cliInstance
}

func NewCli(name string, version string) (*Cli, error) {
	if cliInstance == nil {
		cliInstance, err := createCli(name, version)
		if err != nil {
			return cliInstance, err
		}
		return cliInstance, nil
	}
	return cliInstance, nil
}

func createCli(name string, version string) (*Cli, error) {
	commander := command.GetCommander()
	commander.SetOperator(DEFAULT_OPERATOR)
	cli := &Cli{
		commander:    commander,
		Name:         name,
		HistoryLimit: DEFAULT_HISTORY_LIMIT,
		Symbol:       DEFAULT_SYMBOL,
	}
	err := cli.AddCommand(command.ExitCommand())
	if err != nil {
		return cli, err
	}
	err = cli.AddCommand(command.HelpCommand(""))
	if err != nil {
		return cli, err
	}
	cli, err = cli.SetVersion(version)
	if err != nil {
		return cli, err
	}
	return cli, nil
}

func (cli *Cli) GetVersion() string {
	return command.GetVersionString()
}

func (cli *Cli) SetOperator(operator operator.Operator) *Cli {
	cli.commander.SetOperator(operator)
	return cli
}

func (cli *Cli) SetVersion(version string) (*Cli, error) {
	// Define the regex pattern for the version format vX.Y.Z
	versionRegex := `^\d+\.\d+\.\d+$`
	matched, err := regexp.MatchString(versionRegex, version)
	if err != nil {
		return nil, fmt.Errorf("failed to validate version: %w", err)
	}

	if !matched {
		return nil, errors.New("invalid version format. Expected format: X.Y.Z (e.g., 1.0.0)")
	}

	cli.AddCommand(command.VersionCommand(version))
	return cli, nil
}

func (cli *Cli) SetHelpText(helpText string) *Cli {
	cli.AddCommand(command.HelpCommand(helpText))
	return cli
}

func (cli *Cli) AddCommand(command command.Command) error {
	err := command.Validate()
	if err != nil {
		return err
	}
	cli.commander.AddCommand(command.String(), command)
	return nil
}

func (cli *Cli) Run(interactiveMode bool) {
	args := os.Args
	if len(args) > 1 {
		err := cli.commander.Run(args[1:])
		if err != nil {
			cli.commander.Write(err.Display())
		}
		cli.commander.Write("\n")
	} else if !interactiveMode {
		cli.commander.Write("Interactive shell is disabled!\n")
	} else {
		line, err_ := readline.New(cli.Name + "> ")
		if err_ != nil {
			log.Fatalf("Error initializing readline: %v", err_)
		}
		defer line.Close()
		line.Config.HistoryLimit = cli.HistoryLimit
		for {
			input, err_ := line.Readline()
			if err_ != nil {
				fmt.Println("\nExiting...") // Exit on EOF (Ctrl+D)
				break
			}
			line.SaveHistory(input)
			err := cli.commander.Run(strings.Split(strings.TrimSpace(strings.TrimSuffix(input, "\n")), " "))
			if err != nil {
				cli.commander.Write(err.Display())
			}
			cli.commander.Write("\n")
		}
	}
}
