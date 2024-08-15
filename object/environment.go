// object/environment.go

package object

import "fmt"

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	m := make(map[string]*Module)
	return &Environment{store: s, modules: m, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store   map[string]Object
	modules map[string]*Module
	outer   *Environment
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

func (e *Environment) GetModule(name string) (*Module, bool) {
	module, ok := e.modules[name]
	if !ok && e.outer != nil {
		module, ok = e.outer.GetModule(name)
	}
	return module, ok
}

func (e *Environment) SetModule(name string, module *Module) *Module {
	root := e
	for root.outer != nil {
		root = root.outer
	}
	root.modules[name] = module
	return module
}

func (e *Environment) PrintModules() {
	for k := range e.modules {
		fmt.Printf("mod: %s\n", k)
	}
}

func (e *Environment) Root() *Environment {
	root := e
	for root.outer != nil {
		root = root.outer
	}
	return root
}

type MetaData = map[string]interface{}

type EvalContext struct {
	MetaData       *MetaData
	Line           int
	Column         int
	IsTailPosition bool
}

func NewEvalContext() *EvalContext {
	return &EvalContext{
		MetaData: &MetaData{
			"args": make([]string, 0),
		},
		Line:           1,
		Column:         1,
		IsTailPosition: false,
	}
}
