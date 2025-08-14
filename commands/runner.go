package commands

import (
	"log"
	"slices"
)

type Runner struct {
	Env     string
	Command []string
}

func (ex *Runner) Perform() {
	command := ex.Command[0]

	allCommands := []string{
		"db:create",
		"db:drop",
		"db:migrate",
		"db:rollback",
	}

	if slices.Contains(allCommands, command) {
		log.Println("Running command:", command)
	} else {
		log.Fatal("Running command:", command)
	}
}
