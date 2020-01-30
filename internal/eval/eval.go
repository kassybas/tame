package eval

import (
	"fmt"

	"github.com/antonmedv/expr"
)

type Env map[string]interface{}

func EvaluateExpression(expression string, vars map[string]interface{}) (interface{}, error) {
	env := Env(vars)
	for fnName, fn := range AllFunctions() {
		env[fnName] = fn
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate expression: %s\n\t%s", expression, err.Error())
	}
	return expr.Run(program, env)
}
