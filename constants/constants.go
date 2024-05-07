// constants/constants.go

package constants

import "renelle/object"

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL   = &object.Atom{Value: "nil"}
	OK    = &object.Atom{Value: "ok"}
)
