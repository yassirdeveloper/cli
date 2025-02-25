package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	commands "github.com/yassirdeveloper/cli/commands"
)

func main() {
	commander := commands.NewCommander()

	commander.AddCommand("help", commands.HelpCommand)
	commander.AddCommand("exit", commands.ExitCommand)

	reader := bufio.NewReader(os.Stdin)
	commander.SetWriter(os.Stdout)
	for {
		fmt.Print("> ")
		commandString, _ := reader.ReadString('\n')
		err := commander.Run(strings.TrimSuffix(commandString, "\n"))
		if err != nil {
			commander.Write(err.Display())
		}
		commander.Write("\n")
	}
}
