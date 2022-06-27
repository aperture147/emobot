package eval

import (
	"emobot/bot/application"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const allowedImports = `import ("fmt")`

type commandCollection struct {
	execEnv map[string]*interp.Interpreter
}

func NewCommandCollection() application.CommandCollection {
	execEnv := make(map[string]*interp.Interpreter)
	return commandCollection{execEnv}
}

func (c commandCollection) GetAllCommands() []application.Command {
	initCommand := NewInitEvalCommand(c.execEnv)
	runCommand := NewEvalCommand(c.execEnv)
	return []application.Command{initCommand, runCommand}
}

func newEvalEnvironment() *interp.Interpreter {
	exec := interp.New(interp.Options{})
	_ = exec.Use(stdlib.Symbols)
	_, _ = exec.Eval(allowedImports)
	return exec
}
