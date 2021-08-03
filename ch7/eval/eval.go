// Package eval provides an expression evaluator.
package eval

import (
	"fmt"
	"math"
)

// To evaluate an expression containing variables, we’ll need an environment
// that maps variable names to values.
type Env map[Var]float64

// Performs an environment lookup, which returns zero if the variable is not
// defined.
func (v Var) Eval(env Env) float64 {
	return env[v]
}

// Simply returns the literal value.
func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

// Recursively evaluate their operands, then apply the operation `op` to them.
func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

// Recursively evaluate their operands, then apply the operation `op` to them.
func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		// We don’t consider divisions by zero or infinity to be errors, since
		// they produce a result, albeit non-finite.
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", b.op))
}

// Evaluates the arguments to the `pow`, `sin`, or `sqrt` function, then calls
// the corresponding function in the `math` package.
func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}
