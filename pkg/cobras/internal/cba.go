package internal

type Cba struct {
	cbas    *Cbas
	command *Command
}

func NewCba(cobras *Cbas, command *Command) CbaInterface {
	return &Cba{cbas: cobras, command: command}
}

func (c *Cba) RegisterCommand(name string, cmd CommandInterface) {
	command := NewCommand(name, c.command, c.cbas.commands, cmd)
	c.cbas.commands[command.Index] = command
}
