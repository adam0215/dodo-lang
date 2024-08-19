package object

type Environment struct {
	store    map[string]Object
	mutables map[string]bool
	outer    *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	m := make(map[string]bool)

	return &Environment{store: s, mutables: m, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) IsMutable(name string) bool {
	return e.mutables[name]
}

func (e *Environment) Set(name string, mutable bool, val Object) Object {
	e.store[name] = val

	if mutable {
		e.mutables[name] = true
	}

	return val
}
