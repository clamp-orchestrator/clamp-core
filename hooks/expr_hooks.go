package hooks

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	log "github.com/sirupsen/logrus"
)

// ContextPrefix ... It will contain a prefix so that it will be applicable during condition check
const ContextPrefix = "context."

// ExprHook : Expression hook to validate the condition based branching
type ExprHook struct {
}

// TransformRequest : Transformation will be applied to the given map
func (e *ExprHook) TransformRequest(m map[string]interface{}, s map[string]interface{}) (map[string]interface{}, error) {
	return m, nil
}

// ShouldStepExecute : Check whether Step should execute or skipped
func (e *ExprHook) ShouldStepExecute(
	whenCondition string, stepRequest map[string]interface{}, prefix string) (canStepExecute bool, _ error) {
	log.Debugf("%s Pre-step execution for step is in progress", prefix)

	if !strings.HasPrefix(whenCondition, ContextPrefix) {
		whenCondition = ContextPrefix + whenCondition
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
	}
	return false, fmt.Errorf("invalid boolean expression : %s", whenCondition)
}

// GetExprHook : Getter for Expression hook
func GetExprHook() Hook {
	return &ExprHook{}
}
