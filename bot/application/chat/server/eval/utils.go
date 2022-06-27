package eval

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const allowedImports = ""

func initEvalEnvironment() *interp.Interpreter {
	exec := interp.New(interp.Options{})
	_ = exec.Use(stdlib.Symbols)
}
