package commands

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
)

const OptionLetterPrefix = "-"
const OptionNamePrefix = "--"

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
	String() string
}

type commandInput struct {
	arguments map[string]any
	options   map[string]any
}

func (c *commandInput) String() string {
	return ""
}

func (c *commandInput) ParseArgument(arg commandArgument) (any, Error) {
	argValue := c.arguments[arg.label]
	argValue, err := ParseValue(arg.valueType, argValue)
	if err != nil {
		return nil, &CommandError{message: "Invalid type for argument: " + arg.label}
	}
	return argValue, nil
}

func (c *commandInput) ParseOption(opt commandOption) (any, Error) {
	optValue := c.options[opt.label]
	if optValue == nil {
		return nil, nil
	}
	if opt.valueType == NoType {
		return c.options[opt.label], nil
	}
	optValue, err := ParseValue(opt.valueType, optValue)
	if err != nil {
		return nil, &CommandError{message: "Invalid type for option: " + opt.label}
	}
	return optValue, nil
}

type CommandHanlder func(CommandInput, io.Writer) Error

type Command interface {
	setName(string) Command
	AddArgument(commandArgument) (Command, Error)
	AddOption(commandOption) (Command, Error)
	setHandler(CommandHanlder) Command
	Validate() Error
	Handle(CommandInput, io.Writer) Error
	Parse([]string) (CommandInput, Error)
	String() string
	Help() string
}

type command struct {
	Name        string
	Arguments   []commandArgument
	Options     []commandOption
	handler     CommandHanlder
	Description string
}

func NewCommand(name string, description string, handler CommandHanlder) Command {
	command := &command{
		Name:        name,
		Description: description,
	}
	command.setHandler(handler)
	return command
}

func (c *command) String() string {
	return c.Name
}

func (c *command) Help() string {
	usage := "Usage: > " + c.Name
	if len(c.Arguments) > 0 {
		usage += " [arguments]"
	}
	if len(c.Options) > 0 {
		usage += " [options]"
	}
	return c.Description + " [" + usage + "]"
}

func (c *command) setName(name string) Command {
	c.Name = name
	return c
}

func (c *command) AddArgument(arg commandArgument) (Command, Error) {
	for _, argument := range c.Arguments {
		if argument.label == arg.label {
			return nil, &SetupError{message: fmt.Sprintf("Argument %s for command %s already exists!", arg.label, c.Name)}
		}
	}
	c.Arguments = append(c.Arguments, arg)
	return c, nil
}

func (c *command) AddOption(opt commandOption) (Command, Error) {
	for _, option := range c.Options {
		if option.label == opt.label {
			return nil, &SetupError{message: fmt.Sprintf("Argument %s for command %s already exists!", opt.label, c.Name)}
		}
	}
	c.Options = append(c.Options, opt)
	return c, nil
}

func (c *command) setHandler(commandHandler CommandHanlder) Command {
	c.handler = commandHandler
	return c
}

func (c *command) Validate() Error {
	if c.Name == "" {
		return &SetupError{message: "command name cannot be empty"}
	}
	if len([]rune(c.Name)) < 2 {
		return &SetupError{message: fmt.Sprintf("Command name %s is invalid, needs to be atleast 2 characters long!", c.Name)}
	}
	if len(strings.Split(c.Description, " ")) < 2 {
		return &SetupError{message: fmt.Sprintf("Command %s is invalid, needs to have atleast 2 words long in its description: %s!", c.Name, c.Description)}
	}
	if c.handler == nil {
		return &SetupError{message: fmt.Sprintf("Command %s is not properly set up, needs to have a handler!", c.Name)}
	}
	return nil
}

func (c *command) Handle(input CommandInput, writer io.Writer) Error {
	return c.handler(input, writer)
}

func (c *command) Parse(input []string) (CommandInput, Error) {
	inputLength := len(input)
	inputArgs := make(map[string]any)
	inputOpts := make(map[string]any)

	// Parse arguments
	nbrArguments := len(c.Arguments)
	if inputLength < nbrArguments {
		return nil, &InvalidCommandUsageError{command: c}
	}
	for _, arg := range c.Arguments {
		value, err := ParseValue(arg.valueType, input[arg.position])
		if err != nil {
			return nil, &InvalidCommandUsageError{command: c}
		}
		inputArgs[arg.label] = value
		input = slices.Delete(input, arg.position, arg.position+1)
	}

	// Parse options
	for _, opt := range c.Options {
		index := slices.Index(input, OptionLetterPrefix+string(opt.letter))
		if index == -1 {
			index = slices.Index(input, OptionNamePrefix+opt.name)
		}
		if index != -1 {
			if opt.valueType == NoType {
				inputOpts[opt.label] = true
				input = slices.Delete(input, index, index)
			} else {
				if index+1 >= inputLength {
					return nil, &InvalidCommandUsageError{command: c}
				}
				value, err := ParseValue(opt.valueType, input[index+1])
				if err != nil {
					return nil, &InvalidCommandUsageError{command: c}
				}
				inputOpts[opt.label] = value
				input = slices.Delete(input, index, index+2)
			}
		}
	}

	if len(input) > 0 {
		if strings.HasPrefix(input[0], OptionLetterPrefix) || strings.HasPrefix(input[0], OptionNamePrefix) {
			return nil, &UnreconizedFlagError{command: c.Name, flag: input[0]}
		}
		return nil, &InvalidCommandUsageError{command: c}
	}

	return &commandInput{
		arguments: inputArgs,
		options:   inputOpts,
	}, nil
}

type Commander interface {
	Get(string) (Command, bool)
	AddCommand(string, Command) Commander
	GetCommands() []string
	SetWriter(io.Writer) Commander
	Write(string) Error
	Run([]string) Error
}

type commander struct {
	commands map[string]Command
	writer   io.Writer
}

var commanderInstance Commander

func GetCommander() Commander {
	if commanderInstance == nil {
		commanderInstance = &commander{commands: make(map[string]Command)}
		return commanderInstance
	}
	return commanderInstance
}

func (c *commander) AddCommand(commandName string, command Command) Commander {
	command.setName(commandName)
	c.commands[strings.ToLower(commandName)] = command
	return c
}

func (c *commander) Get(commandName string) (Command, bool) {
	command := c.commands[commandName]
	if command == nil {
		return nil, false
	}
	return command, true
}

func (c *commander) GetCommands() []string {
	return slices.Collect(maps.Keys(c.commands))
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
	command, exists := c.Get(commandName)
	if !exists {
		return &InvalidCommandError{command: commandName}
	}
	input := in[1:]
	inputCommand, err := command.Parse(input)
	if err != nil {
		return err
	}
	return command.Handle(inputCommand, c.writer)
}
