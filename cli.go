package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	readline "github.com/chzyer/readline"
	commands "github.com/yassirdeveloper/cli/commands"
)

const DEFAULT_SYMBOL = ">"
const DEFAULT_HISTORY_LIMIT = 100

var DEFAULT_WRITER = os.Stdout

var cliInstance *Cli

type Cli struct {
	commander    commands.Commander
	Name         string
	HistoryLimit int
	Symbol       string
}

func NewCli(name string) *Cli {
	if cliInstance == nil {
		cliInstance = createCli(name)
		return cliInstance
	}
	return cliInstance
}

func createCli(name string) *Cli {
	commander := commands.GetCommander()
	commander.SetWriter(DEFAULT_WRITER)
	cli := &Cli{
		commander:    commander,
		Name:         name,
		HistoryLimit: DEFAULT_HISTORY_LIMIT,
		Symbol:       DEFAULT_SYMBOL,
	}
	cli.commander.AddCommand("exit", commands.ExitCommand()).
		AddCommand("version", commands.VersionCommand()).
		AddCommand("help", commands.HelpCommand())
	return cli
}

func (cli *Cli) SetWriter(writer io.Writer) *Cli {
	cli.commander.SetWriter(writer)
	return cli
}

func (cli *Cli) AddCommand(command commands.Command) error {
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
