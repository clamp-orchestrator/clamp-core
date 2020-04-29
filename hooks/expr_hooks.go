package hooks

import (
	"fmt"
	"github.com/antonmedv/expr"
	"log"
	"strings"
)

type ExprHook struct {
}

func (e *ExprHook) TransformRequest(m map[string]interface{}, s map[string]interface{}) (map[string]interface{}, error) {
	return m, nil
}

func (e *ExprHook) ShouldStepExecute(whenCondition string, stepRequest map[string]interface{}, prefix string) (canStepExecute bool, _ error) {
	log.Printf("%s Pre-step execution for step is in progress", prefix)
	if !strings.HasPrefix(whenCondition, "context.") {
		whenCondition = "content." + whenCondition
	}
	env := map[string]interface{}{
		"context": stepRequest,
	}
	// Compile code into bytecode. This step can be done once and program may be reused.
	// Specify environment for type check.
	program, err := expr.Compile(whenCondition, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	if canStepExecute, ok := output.(bool); ok {
		return canStepExecute, nil
	} else {
		return false, fmt.Errorf("invalid boolean expression : %s", whenCondition)
	}
}

func GetExprHook() Hook {
	return &ExprHook{}
}
