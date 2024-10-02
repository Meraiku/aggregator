package cli

import "fmt"

func Login(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return ErrNoArgs
	}

	err := state.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been set!\n", state.cfg.CurrentUserName)

	return nil
}
