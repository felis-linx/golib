package utils

import (
	"reflect"
	"strconv"
	"strings"
)

type encoderFunc func(reflect.Value) string

func encodeBool(v reflect.Value) string {
	return strconv.FormatBool(v.Bool())
}

func encodeInt(v reflect.Value) string {
	return strconv.FormatInt(int64(v.Int()), 10)
}

func encodeUint(v reflect.Value) string {
	return strconv.FormatUint(uint64(v.Uint()), 10)
}

func encodeFloat(v reflect.Value, bits int) string {
	return strconv.FormatFloat(v.Float(), 'f', 6, bits)
}

func encodeFloat32(v reflect.Value) string {
	return encodeFloat(v, 32)
}

func encodeFloat64(v reflect.Value) string {
	return encodeFloat(v, 64)
}

func encodeString(v reflect.Value) string {
	return v.String()
}

func typeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {

	case reflect.Bool:
		return encodeBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeUint

	case reflect.Float32:
		return encodeFloat32

	case reflect.Float64:
		return encodeFloat64

	case reflect.Ptr:
		f := typeEncoder(t.Elem())
		return func(v reflect.Value) string {
			if v.IsNil() {
				return "null"
			}
			return f(v.Elem())
		}
	case reflect.String:
		return encodeString
	default:
		return nil
	}
}

func isValidStructPointer(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Ptr && v.Elem().IsValid() && v.Elem().Type().Kind() == reflect.Struct
}

func isValidInterface(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Interface && v.Elem().IsValid()
}

type tagOptions []string

func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func:
	case reflect.Map, reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		type zero interface {
			IsZero() bool
		}
		if v.Type().Implements(reflect.TypeOf((*zero)(nil)).Elem()) {
			iz := v.MethodByName("IsZero").Call([]reflect.Value{})[0]
			return iz.Interface().(bool)
		}
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}

	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

func mask(tag string, v reflect.Value) (result map[string]interface{}) {
	t := v.Type()
	result = map[string]interface{}{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldTag := field.Tag.Get(tag)
		name, opts := parseTag(fieldTag)

		if name == "" || opts.Contains("skip") {
			continue
		}

		if opts.Contains("mask") {
			result[name] = "***"
			continue
		}

		if isValidStructPointer(v.Field(i)) {
			result[name] = mask(tag, v.Field(i).Elem())
			continue
		}

		encFunc := typeEncoder(v.Field(i).Type())
		if encFunc != nil {
			value := encFunc(v.Field(i))
			if opts.Contains("omitempty") && isZero(v.Field(i)) {
				continue
			}

			result[name] = value
			continue
		}

		if v.Field(i).Type().Kind() == reflect.Struct {
			result[name] = mask(tag, v.Field(i))
			continue
		}

		if isValidInterface(v.Field(i)) {
			result[name] = mask(tag, reflect.ValueOf(v.Field(i).Interface()))
			continue
		}
	}

	return
}

// MaskValues (tag string, data interface{}) (masked map[string]interface{}, err error)
// Mask value with 'mask' option in metadata tag,
// and exclude fields marked with 'skip' option
func MaskValues(tag string, data interface{}) (masked map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(*reflect.ValueError)
		}
	}()

	masked = mask(tag, reflect.ValueOf(data))

	return
}
