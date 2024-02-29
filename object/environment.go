// object/environment.go

package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

type MetaData = map[string]interface{}

type EvalContext struct {
	MetaData *MetaData
	Line     int
	Column   int
}

func NewEvalContext() *EvalContext {
	return &EvalContext{
		MetaData: &MetaData{
			"args": make([]string, 0),
		},
		Line:   1,
		Column: 1,
	}
}
