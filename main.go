package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	readline "github.com/chzyer/readline"
	commands "github.com/yassirdeveloper/cli/commands"
)

const HistoryLimit = 100
const AppName = "cli"

func main() {
	commander := commands.GetCommander()
	commander.AddCommand("exit", commands.ExitCommand())
	commander.AddCommand("version", commands.VersionCommand())
	commander.AddCommand("help", commands.HelpCommand())
	commander.SetWriter(os.Stdout)

	args := os.Args
	if len(args) > 1 {
		err := commander.Run(args[1:])
		if err != nil {
			commander.Write(err.Display())
			os.Exit(1)
		}
		commander.Write("\n")
		os.Exit(0)
	}

	line, err_ := readline.New(AppName + "> ")
	if err_ != nil {
		log.Fatalf("Error initializing readline: %v", err_)
	}
	defer line.Close()
	line.Config.HistoryLimit = HistoryLimit

	for {
		input, err_ := line.Readline()
		if err_ != nil {
			fmt.Println("\nExiting...") // Exit on EOF (Ctrl+D)
			break
		}
		line.SaveHistory(input)
		err := commander.Run(strings.Split(strings.TrimSpace(strings.TrimSuffix(input, "\n")), " "))
		if err != nil {
			commander.Write(err.Display())
		}
		commander.Write("\n")
	}
}
