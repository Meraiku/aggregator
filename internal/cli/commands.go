package cli

import "sync"

type CLIFunc func(*State, Command) error

type Command struct {
	Name string
	Args []string
}

func NewCommand(name string, args []string) *Command {
	return &Command{
		Name: name,
		Args: args,
	}
}

type Commands struct {
	cmds map[string]CLIFunc
	mu   sync.Mutex
}

func NewCommands() *Commands {
	return &Commands{
		cmds: make(map[string]CLIFunc),
		mu:   sync.Mutex{},
	}
}

func (c *Commands) Register(name string, f CLIFunc) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.cmds[name]; ok {
		return ErrCommandAlreadyRegisterd
	}

	c.cmds[name] = f

	return nil
}

func (c *Commands) RegisterHandlers() error {
	funcs := []struct {
		name string
		f    CLIFunc
	}{
		{
			name: "login",
			f:    Login,
		},
		{
			name: "register",
			f:    Register,
		},
		{
			name: "reset",
			f:    Reset,
		},
		{
			name: "users",
			f:    Users,
		},
		{
			name: "agg",
			f:    Agg,
		},
	}

	for _, f := range funcs {
		if err := c.Register(f.name, f.f); err != nil {
			return err
		}
	}

	return nil
}

func (c *Commands) Run(state *State, cmd Command) error {
	if state == nil {
		return ErrUnknownState
	}

	f, ok := c.cmds[cmd.Name]
	if !ok {
		return ErrUnknownCommand
	}

	return f(state, cmd)
}
