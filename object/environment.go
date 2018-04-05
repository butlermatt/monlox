package object

type Environment struct {
	store map[string]Object
}

// NewEnvironment returns an environment to scope.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get retrieves an variable value stored in the Environment. Returns the Object plus boolean if it was successful or not.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set stores a variable and value to the Environment. Returns the stored value.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
