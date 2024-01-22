package validate

import (
	"reflect"
	"regexp"
	"strconv"
)

var Required = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	switch memberVal.Kind() {
	case reflect.String:
		return memberVal.String() != ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return memberVal.Len() > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return memberVal.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return memberVal.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return memberVal.Float() != 0
	case reflect.Ptr, reflect.Interface:
		return !memberVal.IsNil()
	default:
		return false
	}
}

var Min = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	i, err := strconv.Atoi(condition)
	if err != nil {
		panic(err)
	}
	switch memberVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return memberVal.Int() >= int64(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return memberVal.Uint() >= uint64(i)
	case reflect.Float32, reflect.Float64:
		return memberVal.Float() >= float64(i)
	default:
		return false
	}
}

var Max = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	i, err := strconv.Atoi(condition)
	if err != nil {
		panic(err)
	}
	switch memberVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return memberVal.Int() <= int64(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return memberVal.Uint() <= uint64(i)
	case reflect.Float32, reflect.Float64:
		return memberVal.Float() <= float64(i)
	default:
		return false
	}
}

var Len = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	i, err := strconv.Atoi(condition)
	if err != nil {
		panic(err)
	}
	switch memberVal.Kind() {
	case reflect.String:
		return len(memberVal.String()) == i
	case reflect.Slice, reflect.Array, reflect.Map:
		return memberVal.Len() == i
	default:
		return false
	}
}

var MinLen = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	i, err := strconv.Atoi(condition)
	if err != nil {
		panic(err)
	}
	switch memberVal.Kind() {
	case reflect.String:
		return len(memberVal.String()) >= i
	case reflect.Slice, reflect.Array, reflect.Map:
		return memberVal.Len() >= i
	default:
		return false
	}
}

var MaxLen = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	i, err := strconv.Atoi(condition)
	if err != nil {
		panic(err)
	}
	switch memberVal.Kind() {
	case reflect.String:
		return len(memberVal.String()) <= i
	case reflect.Slice, reflect.Array, reflect.Map:
		return memberVal.Len() <= i
	default:
		return false
	}
}

var Email = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	email := memberVal.String()
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

var Regex = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(condition)
	return regex.MatchString(memberVal.String())
}

var Alpha = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNum = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaDash = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return regex.MatchString(memberVal.String())
}

var Numeric = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[0-9]+$`)
	return regex.MatchString(memberVal.String())
}

var NumericDash = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[0-9_-]+$`)
	return regex.MatchString(memberVal.String())
}

var NumericDot = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[0-9.]+$`)
	return regex.MatchString(memberVal.String())
}

var NumericComma = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[0-9,]+$`)
	return regex.MatchString(memberVal.String())
}

var NumericDashDot = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[0-9_.-]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumeric = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericDash = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericDot = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9.]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericComma = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9,]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericDashDot = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaSpace = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaDashSpace = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z_\- ]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericSpace = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericDashSpace = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_\- ]+$`)
	return regex.MatchString(memberVal.String())
}

var AlphaNumericDashSpaceDot = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_.\- ]+$`)
	return regex.MatchString(memberVal.String())
}

var IP = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
	return regex.MatchString(memberVal.String())
}

var IPV4 = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
	return regex.MatchString(memberVal.String())
}

var IPV6 = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	return regex.MatchString(memberVal.String())
}

var Base64 = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9+/]+={0,2}$`)
	return regex.MatchString(memberVal.String())
}

var Base64URL = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9-_]+={0,2}$`)
	return regex.MatchString(memberVal.String())
}

var Base64RawURL = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._-]+={0,2}$`)
	return regex.MatchString(memberVal.String())
}

var Hexadecimal = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^[a-fA-F0-9]+$`)
	return regex.MatchString(memberVal.String())
}

var HexColor = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^#?([a-fA-F0-9]{3}|[a-fA-F0-9]{6})$`)
	return regex.MatchString(memberVal.String())
}

var RGBColor = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^rgb\((\d{1,3}),(\d{1,3}),(\d{1,3})\)$`)
	return regex.MatchString(memberVal.String())
}

var RGBAColor = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^rgba\((\d{1,3}),(\d{1,3}),(\d{1,3}),([01]|0\.\d+)\)$`)
	return regex.MatchString(memberVal.String())
}

var HSLColor = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^hsl\((\d{1,3}),(\d{1,3})%,(\d{1,3})%\)$`)
	return regex.MatchString(memberVal.String())
}

var HSLAColor = func(memberVal reflect.Value, condition string, structVal reflect.Value) bool {
	regex := regexp.MustCompile(`^hsla\((\d{1,3}),(\d{1,3})%,(\d{1,3})%,([01]|0\.\d+)\)$`)
	return regex.MatchString(memberVal.String())
}
