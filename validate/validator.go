package validate

import (
	"errors"
	"github.com/samber/lo"
	"reflect"
	"strings"
)

type Validator struct {
	Checkers            map[string]Checker
	StructErrDefaultMsg string
	StructValidTag      string
	StructErrMsgTag     string
	DefaultErrorMsg     string
}

func New() *Validator {
	return &Validator{
		Checkers: map[string]Checker{
			"required":                 Required,
			"min":                      Min,
			"max":                      Max,
			"len":                      Len,
			"minLen":                   MinLen,
			"maxLen":                   MaxLen,
			"email":                    Email,
			"regex":                    Regex,
			"alpha":                    Alpha,
			"alphaNum":                 AlphaNum,
			"alphaDash":                AlphaDash,
			"numeric":                  Numeric,
			"numericDash":              NumericDash,
			"numericDot":               NumericDot,
			"numericComma":             NumericComma,
			"numericDashDot":           NumericDashDot,
			"alphaNumeric":             AlphaNumeric,
			"alphaNumericDash":         AlphaNumericDash,
			"alphaNumericDot":          AlphaNumericDot,
			"alphaNumericComma":        AlphaNumericComma,
			"alphaNumericDashDot":      AlphaNumericDashDot,
			"alphaSpace":               AlphaSpace,
			"alphaDashSpace":           AlphaDashSpace,
			"alphaNumericSpace":        AlphaNumericSpace,
			"alphaNumericDashSpace":    AlphaNumericDashSpace,
			"alphaNumericDashSpaceDot": AlphaNumericDashSpaceDot,
			"IP":                       IP,
			"IPv4":                     IPV4,
			"IPv6":                     IPV6,
			"base64":                   Base64,
			"Base64URL":                Base64URL,
			"hexadecimal":              Hexadecimal,
			"hexColor":                 HexColor,
			"RGBColor":                 RGBColor,
			"RGBAColor":                RGBAColor,
			"HSLColor":                 HSLColor,
			"HSLAColor":                HSLAColor,
			"oneOf":                    OneOf,
			"sameAs":                   SameAs,
			"cnPhone":                  CnPhone,
		},
		StructValidTag:      "validate",
		StructErrDefaultMsg: "invalid param.",
		StructErrMsgTag:     "msg",
	}
}

func (v *Validator) AddCheckers(cs map[string]Checker) {
	for s, checker := range cs {
		v.Checkers[s] = checker
	}
}

func (v *Validator) CheckStruct(s any) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			pv := val.Field(i)
			field := val.Type().Field(i)
			value, ok := field.Tag.Lookup(v.StructValidTag)
			if !ok || value == "" {
				continue
			}
			if v.doCheck(pv, value, val) == false {
				return errors.New(v.fetchErrorMsg(field))
			}
		}

		return nil
	}

	return nil
}

func (v *Validator) CheckVar(s any, validates string) bool {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return v.doCheck(val, validates, val)
}

func (v *Validator) doCheck(pv reflect.Value, validate string, sv reflect.Value) bool {
	cna := strings.Split(validate, ",")
	for i, v := range cna {
		cna[i] = strings.TrimSpace(v)
	}

	if lo.IndexOf(cna, "required") > 0 &&
		v.Checkers["required"](pv, "", sv) != true {
		return false
	}

	for _, checkerName := range cna {
		if checkerName == "" || checkerName == "required" {
			continue
		}

		kv := strings.Split(checkerName, "=")
		checkerName = strings.TrimSpace(kv[0])
		if len(kv) == 1 {
			if cn, ok := v.Checkers[checkerName]; ok {
				if cn(pv, "", sv) == false {
					return false
				}
			}
			continue
		}

		param := strings.TrimSpace(kv[1])
		if cn, ok := v.Checkers[checkerName]; ok {
			if cn(pv, param, sv) == false {
				return false
			}
		}
	}

	return true
}

func (v *Validator) fetchErrorMsg(field reflect.StructField) string {
	value, ok := field.Tag.Lookup(v.StructErrMsgTag)
	if !ok || value == "" {
		return v.DefaultErrorMsg
	}
	return value
}
