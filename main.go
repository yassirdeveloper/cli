package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	readline "github.com/chzyer/readline"
	commands "github.com/yassirdeveloper/cli/commands"
)

func main() {
	commander := commands.NewCommander()

	commander.AddCommand("help", commands.HelpCommand)
	commander.AddCommand("exit", commands.ExitCommand)
	commander.SetWriter(os.Stdout)
	args := os.Args
	if len(args) > 1 {
		err := commander.Run(args[1:])
		if err != nil {
			commander.Write(err.Display())
			os.Exit(1)
		}
		os.Exit(0)
	}
	line, err_ := readline.New("> ")
	if err_ != nil {
		log.Fatalf("Error initializing readline: %v", err_)
	}
	defer line.Close()
	line.Config.HistoryLimit = 100
	for {
		input, err_ := line.Readline()
		if err_ != nil {
			fmt.Println("\nExiting...") // Exit on EOF (Ctrl+D)
			break
		}
		input = strings.TrimSpace(input)
		line.SaveHistory(input)

		err := commander.Run(strings.Split(strings.TrimSpace(strings.TrimSuffix(input, "\n")), " "))
		if err != nil {
			commander.Write(err.Display())
		}
		commander.Write("\n")

	}
}
