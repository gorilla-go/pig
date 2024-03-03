package injector

import (
	"github.com/gorilla-go/pig/param"
	"github.com/samber/lo"
	"reflect"
	"unsafe"
)

func ServiceInjector(tp reflect.Type, t any, at unsafe.Pointer) {
	reflect.NewAt(tp, at).Elem().Set(reflect.ValueOf(t))
}

func RequestInjector(tp reflect.Type, val *param.RequestParamItems[string], at unsafe.Pointer) {
	if val != nil && len(val.GetParams()) > 0 && CanInjected(tp.Kind()) {
		reflect.NewAt(tp, at).Elem().Set(
			reflect.ValueOf(ConvertStringToKind(val, tp)),
		)
	}
}

func CanInjected(k reflect.Kind) bool {
	return lo.IndexOf([]reflect.Kind{
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Slice,
		reflect.String,
	}, k) != -1
}

func ConvertStringToKind(s *param.RequestParamItems[string], k reflect.Type) any {
	switch k.Kind() {
	case reflect.Bool:
		return s.Bool()
	case reflect.Int:
		return s.Int()
	case reflect.Int8:
		return int8(s.Int())
	case reflect.Int16:
		return int16(s.Int())
	case reflect.Int32:
		return int32(s.Int())
	case reflect.Int64:
		return s.Int64()
	case reflect.Uint:
		return uint(s.Int())
	case reflect.Uint8:
		return uint8(s.Int())
	case reflect.Uint16:
		return uint16(s.Int())
	case reflect.Uint32:
		return uint32(s.Int())
	case reflect.Uint64:
		return uint64(s.Int())
	case reflect.Float32:
		return float32(s.Float64())
	case reflect.Float64:
		return s.Float64()
	case reflect.Slice:
		sType := reflect.SliceOf(k.Elem())
		l := len(s.GetParams())
		slice := reflect.MakeSlice(sType, l, l)
		for i := 0; i < l; i++ {
			itemType := slice.Index(i).Type()
			if !CanInjected(slice.Index(i).Type().Kind()) {
				panic("unsupported inject type: []" + itemType.String())
			}
			slice.Index(i).Set(reflect.ValueOf(
				ConvertStringToKind(
					param.NewRequestParamItems([]string{s.GetParams()[i].String()}),
					itemType,
				),
			))
		}
		return slice.Interface()
	case reflect.String:
		return s.String()
	}

	return nil
}
