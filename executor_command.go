package main

type CommandExecutor struct {
	commands []string
	Executor
}

func NewCommandExecutor(commands []string) *CommandExecutor {
	e := new(CommandExecutor)
	e.commands = commands

	return e
}

func (e *CommandExecutor) Exec() string {
	for _, cmd := range e.commands {
		_, err := executeCommandWithOutput(cmd)
		if err != nil {
			return StatusFail
		}
	}
	return StatusSuccess
}
