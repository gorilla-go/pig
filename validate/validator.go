package validate

import (
	"errors"
	"reflect"
	"strings"
)

type Checker func(memberVal reflect.Value, condition string, structVal reflect.Value) bool

type IValidator interface {
	Validate(any) error
}

type Validator struct {
	Checkers   map[string]Checker
	DefaultMsg string
}

func New(m map[string]Checker) *Validator {
	return &Validator{
		Checkers:   m,
		DefaultMsg: "The parameter failed to pass validation.",
	}
}

func (v *Validator) Validate(s any) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("validate target must be struct or struct pointer.")
	}

	for i := 0; i < val.NumField(); i++ {
		pv := val.Field(i)
		field := val.Type().Field(i)
		value, ok := field.Tag.Lookup("validate")
		if !ok || value == "" {
			continue
		}

		cna := strings.Split(value, ",")
		for _, checkerName := range cna {
			checkerName = strings.TrimSpace(checkerName)
			if checkerName == "" {
				continue
			}

			kv := strings.Split(checkerName, "=")
			checkerName = kv[0]
			if len(kv) == 1 {
				if cn, ok := v.Checkers[checkerName]; ok {
					if cn(pv, "", val) != true {
						return errors.New(v.fetchErrorMsg(field))
					}
				}
				continue
			}

			param := kv[1]
			if cn, ok := v.Checkers[checkerName]; ok {
				if cn(pv, param, val) != true {
					return errors.New(v.fetchErrorMsg(field))
				}
			}
		}
	}
	return nil
}

func (v *Validator) fetchErrorMsg(field reflect.StructField) string {
	value, ok := field.Tag.Lookup("msg")
	if !ok || value == "" {
		return v.DefaultMsg
	}
	return value
}
