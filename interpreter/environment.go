package interpreter

// Environment stores variable bindings for the interpreter.
type Environment struct {
	store map[string]Value
}

// NewEnvironment creates a new empty environment.
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Value)}
}

// Set binds a value to a variable name.
func (e *Environment) Set(name string, value Value) {
	e.store[name] = value
}

// Get retrieves the value bound to a variable name.
func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.store[name]
	return val, ok
}

// Exists checks if a variable is bound in the environment.
func (e *Environment) Exists(name string) bool {
	_, ok := e.store[name]
	return ok
}
