package command

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/yassirdeveloper/cli/errors"
	"github.com/yassirdeveloper/cli/operator"
)

const OptionLetterPrefix = "-"
const OptionNamePrefix = "--"

type ValueType int

type Stringer interface {
	String() string
}

type CommandArgument struct {
	Label       string
	Description string
	Position    int
	ValueType   ValueType
}

type CommandOption struct {
	Label       string
	Description string
	Letter      rune
	Name        string
	ValueType   ValueType
}

type CommandInput interface {
	ParseArgument(CommandArgument) (any, errors.Error)
	ParseOption(CommandOption) (any, errors.Error)
	String() string
}

type commandInput struct {
	arguments map[string]any
	options   map[string]any
}

func (c *commandInput) String() string {
	return ""
}

func (c *commandInput) ParseArgument(arg CommandArgument) (any, errors.Error) {
	argValue := c.arguments[arg.Label]
	argValue, err := ParseValue(arg.ValueType, argValue)
	if err != nil {
		return nil, &CommandError{message: "Invalid type for argument: " + arg.Label}
	}
	return argValue, nil
}

func (c *commandInput) ParseOption(opt CommandOption) (any, errors.Error) {
	optValue := c.options[opt.Label]
	if optValue == nil {
		return nil, nil
	}
	if opt.ValueType == NoType {
		return c.options[opt.Label], nil
	}
	optValue, err := ParseValue(opt.ValueType, optValue)
	if err != nil {
		return nil, &CommandError{message: "Invalid type for option: " + opt.Label}
	}
	return optValue, nil
}

type CommandHanlder func(CommandInput, operator.Operator) errors.Error

type Command interface {
	setName(string) Command
	AddArgument(CommandArgument) (Command, errors.Error)
	AddOption(CommandOption) (Command, errors.Error)
	setHandler(CommandHanlder) Command
	Validate() errors.Error
	Handle(CommandInput, operator.Operator) errors.Error
	Parse([]string) (CommandInput, errors.Error)
	String() string
	Help() string
}

type command struct {
	Name        string
	Arguments   []CommandArgument
	Options     []CommandOption
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
	usageBuilder := &strings.Builder{}
	usageBuilder.WriteString("Usage: > " + c.Name)
	if len(c.Arguments) > 0 {
		for _, arg := range c.Arguments {
			usageBuilder.WriteString(" " + arg.Label)
		}
	}
	if len(c.Options) > 0 {
		usageBuilder.WriteString(" [options]")
	}

	helpText := fmt.Sprintf("\t- %-15s %s\n", c.Name+":", c.Description+" ["+usageBuilder.String()+"]")
	if len(c.Options) > 0 {
		optionsBuilder := &strings.Builder{}
		for _, opt := range c.Options {
			optionsBuilder.WriteString(fmt.Sprintf("\t   -%c | --%s:  %s.\n", opt.Letter, opt.Label, opt.Description))
		}
		helpText += optionsBuilder.String()
	}
	return helpText
}

func (c *command) setName(name string) Command {
	c.Name = name
	return c
}

func (c *command) AddArgument(arg CommandArgument) (Command, errors.Error) {
	for _, argument := range c.Arguments {
		if argument.Label == arg.Label {
			return nil, errors.NewSetupError(fmt.Sprintf("Argument %s for command %s already exists!", arg.Label, c.Name))
		}
	}
	c.Arguments = append(c.Arguments, arg)
	return c, nil
}

func (c *command) AddOption(opt CommandOption) (Command, errors.Error) {
	for _, option := range c.Options {
		if option.Label == opt.Label {
			return nil, errors.NewSetupError(fmt.Sprintf("Argument %s for command %s already exists!", opt.Label, c.Name))
		}
	}
	c.Options = append(c.Options, opt)
	return c, nil
}

func (c *command) setHandler(commandHandler CommandHanlder) Command {
	c.handler = commandHandler
	return c
}

func (c *command) Validate() errors.Error {
	if c.Name == "" {
		return errors.NewSetupError("command name cannot be empty")
	}
	if len([]rune(c.Name)) < 2 {
		return errors.NewSetupError(fmt.Sprintf("Command name %s is invalid, needs to be atleast 2 characters long!", c.Name))
	}
	if len(strings.Split(c.Description, " ")) < 2 {
		return errors.NewSetupError(fmt.Sprintf("Command %s is invalid, needs to have atleast 2 words long in its description: %s!", c.Name, c.Description))
	}
	if c.handler == nil {
		return errors.NewSetupError(fmt.Sprintf("Command %s is not properly set up, needs to have a handler!", c.Name))
	}
	return nil
}

func (c *command) Handle(input CommandInput, operator operator.Operator) errors.Error {
	return c.handler(input, operator)
}

func (c *command) Parse(input []string) (CommandInput, errors.Error) {
	inputLength := len(input)
	inputArgs := make(map[string]any)
	inputOpts := make(map[string]any)

	// Parse arguments
	nbrArguments := len(c.Arguments)
	if inputLength < nbrArguments {
		return nil, &InvalidCommandUsageError{command: c}
	}
	for _, arg := range c.Arguments {
		value, err := ParseValue(arg.ValueType, input[arg.Position])
		if err != nil {
			return nil, &InvalidCommandUsageError{command: c}
		}
		inputArgs[arg.Label] = value
		input = slices.Delete(input, arg.Position, arg.Position+1)
	}

	// Parse options
	for _, opt := range c.Options {
		index := slices.Index(input, OptionLetterPrefix+string(opt.Letter))
		if index == -1 {
			index = slices.Index(input, OptionNamePrefix+opt.Name)
		}
		if index != -1 {
			if opt.ValueType == NoType {
				inputOpts[opt.Label] = true
				input = slices.Delete(input, index, index)
			} else {
				if index+1 >= inputLength {
					return nil, &InvalidCommandUsageError{command: c}
				}
				value, err := ParseValue(opt.ValueType, input[index+1])
				if err != nil {
					return nil, &InvalidCommandUsageError{command: c}
				}
				inputOpts[opt.Label] = value
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
	SetOperator(operator.Operator) Commander
	Write(string) errors.Error
	Run([]string) errors.Error
}

type commander struct {
	commands map[string]Command
	operator operator.Operator
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

func (c *commander) SetOperator(operator operator.Operator) Commander {
	c.operator = operator
	return c
}

func (c *commander) Write(output string) errors.Error {
	err := c.operator.Write(output)
	if err != nil {
		return errors.NewUnexpectedError(err)
	}
	return nil
}

func (c *commander) Run(in []string) errors.Error {
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
	return command.Handle(inputCommand, c.operator)
}
