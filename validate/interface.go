package validate

import "reflect"

type Checker func(memberVal reflect.Value, condition string, structVal reflect.Value) bool
type VarChecker func(memberVal reflect.Value, condition string) bool

type IValidator interface {
	CheckStruct(structVar any) error
	CheckVar(param any, ValidStr string) bool
}
