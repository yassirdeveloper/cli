package commands

import (
	"io"
	"slices"
	"strings"
)

type ValueType int

type Stringer interface {
	String() string
}

type commandArgument struct {
	label       string
	description string
	position    int
	valueType   ValueType
}

type commandOption struct {
	label       string
	description string
	letter      rune
	name        string
	valueType   ValueType
}

type CommandInput interface {
	ParseArgument(commandArgument) (any, Error)
	ParseOption(commandOption) (any, Error)
}

type commandInput struct {
	arguments map[commandArgument]any
	options   map[commandOption]any
}

func (c *commandInput) ParseArgument(arg commandArgument) (any, Error) {
	argValue := c.arguments[arg]
	argValue, err := ParseValue(arg.valueType, argValue)
	if err != nil {
		return nil, &CommandError{message: "Invalid type for argument: " + arg.label}
	}
	return argValue, nil
}

func (c *commandInput) ParseOption(opt commandOption) (any, Error) {
	optValue := c.options[opt]
	if optValue == nil {
		return nil, nil
	}
	if opt.valueType == NoType {
		return c.options[opt], nil
	}
	optValue, err := ParseValue(opt.valueType, optValue)
	if err != nil {
		return nil, &CommandError{message: "Invalid type for option: " + opt.label}
	}
	return optValue, nil
}

type CommandHanlder func(CommandInput, io.Writer) Error

type Command interface {
	AddArgument(commandArgument) Command
	AddOption(commandOption) Command
	SetHandler(CommandHanlder) Command
	Handle(CommandInput, io.Writer) Error
	Parse([]string) (CommandInput, Error)
	String() string
}

type command struct {
	Name      string
	Arguments []commandArgument
	Options   []commandOption
	handler   CommandHanlder
}

func NewCommand(name string) Command {
	return &command{Name: name}
}

func (c *command) String() string {
	return c.Name
}

func (c *command) AddArgument(arg commandArgument) Command {
	c.Arguments = append(c.Arguments, arg)
	return c
}

func (c *command) AddOption(opt commandOption) Command {
	c.Options = append(c.Options, opt)
	return c
}

func (c *command) SetHandler(commandHandler CommandHanlder) Command {
	c.handler = commandHandler
	return c
}

func (c *command) Handle(input CommandInput, writer io.Writer) Error {
	return c.handler(input, writer)
}

func (c *command) Parse(input []string) (CommandInput, Error) {
	inputLength := len(input)
	inputArgs := make(map[commandArgument]any)
	inputOpts := make(map[commandOption]any)

	// Parse arguments
	nbrArguments := len(c.Arguments)
	if inputLength < nbrArguments {
		return nil, &InvalidCommandUsageError{command: c.Name}
	}
	for _, arg := range c.Arguments {
		value, err := ParseValue(arg.valueType, input[arg.position])
		if err != nil {
			return nil, &InvalidCommandUsageError{command: c.Name}
		}
		inputArgs[arg] = value
	}

	// Parse options
	for _, opt := range c.Options {
		index := slices.Index(input, string(opt.letter))
		if index == -1 {
			index = slices.Index(input, opt.name)
		}
		if index != -1 {
			if opt.valueType == NoType {
				inputOpts[opt] = true
			} else {
				if index+1 >= inputLength {
					return nil, &InvalidCommandUsageError{command: c.Name}
				}
				value, err := ParseValue(opt.valueType, input[index+1])
				if err != nil {
					return nil, &InvalidCommandUsageError{command: c.Name}
				}
				inputOpts[opt] = value
			}
		}
	}

	return &commandInput{
		arguments: inputArgs,
		options:   inputOpts,
	}, nil
}

type Commander interface {
	Get(string) (Command, Error)
	AddCommand(string, Command) Commander
	SetWriter(io.Writer) Commander
	Write(string) Error
	Run([]string) Error
}

type commander struct {
	commands map[string]Command
	writer   io.Writer
}

func NewCommander() Commander {
	return &commander{commands: make(map[string]Command)}
}

func (c *commander) AddCommand(commandName string, command Command) Commander {
	c.commands[strings.ToLower(commandName)] = command
	return c
}

func (c *commander) Get(commandName string) (Command, Error) {
	command := c.commands[commandName]
	if command == nil {
		return nil, &InvalidCommandError{command: commandName}
	}
	return command, nil
}

func (c *commander) SetWriter(writer io.Writer) Commander {
	c.writer = writer
	return c
}

func (c *commander) Write(output string) Error {
	_, err := c.writer.Write([]byte(output))
	if err != nil {
		return &UnexpectedError{err: err}
	}
	return nil
}

func (c *commander) Run(in []string) Error {
	commandName := strings.ToLower(in[0])
	command, err := c.Get(commandName)
	if err != nil {
		return err
	}
	inputCommand, err := command.Parse(in[1:])
	if err != nil {
		return err
	}
	return command.Handle(inputCommand, c.writer)
}
