package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	commander := NewCommander()
	reader := bufio.NewReader(os.Stdin)
	commander.SetWriter(os.Stdout)
	for {
		fmt.Print("> ")
		commandString, _ := reader.ReadString('\n')
		err := commander.Run(commandString)
		if err != nil {
			commander.Write(err.Display())
		}
	}
}
