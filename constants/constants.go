// constants/constants.go

package constants

import "renelle/object"

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL   = &object.Atom{Value: "nil"}
	OK    = &object.Atom{Value: "ok"}
	ERROR = &object.Atom{Value: "error"}
	CONT  = &object.Atom{Value: "cont"}
	HALT  = &object.Atom{Value: "halt"}
	SOME  = &object.Atom{Value: "some"}
	NONE  = &object.Atom{Value: "none"}
)
