package cobras

type CommandChild struct {
	runner  *CommandRunner
	command *Command
}

func NewCommandChild(runner *CommandRunner, command *Command) CommandChildInterface {
	return &CommandChild{runner: runner, command: command}
}

func (c *CommandChild) RegisterCommand(name string, cmd CommandInterface) {
	command := NewCommand(name, NewCommandOptions{
		Parent:   c.command,
		Commands: c.runner.commands,
		Cmd:      cmd,
	})
	c.runner.commands[command.Index] = command
}
