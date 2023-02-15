package structer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

var (
	errInvalidType = errors.New("not a struct pointer")
)

type Structer interface {
	Set(ptr interface{}) error
	MustSet(ptr interface{})
	CanUpdate(v interface{}) bool
	PrettyPrint(v interface{})
}

// SetDefaults initializes members in a struct referenced by a pointer.
// Maps and slices are initialized by `make` and other primitive types are set with default values.
// `ptr` should be a struct pointer
func SetDefaults(ptr interface{}) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return errInvalidType
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return errInvalidType
	}

	for i := 0; i < t.NumField(); i++ {
		if err := setField(v.Field(i), getValue(t.Field(i))); err != nil {
			return err
		}
	}
	return nil
}

// MustSet function is a wrapper of Set function
// It will call Set and panic if err not equals nil.
func MustSet(ptr interface{}) {
	if err := SetDefaults(ptr); err != nil {
		panic(err)
	}
}

// CanUpdate returns true when the given value is an initial value of its type
func CanUpdate(v interface{}) bool {
	return isInitialValue(reflect.ValueOf(v))
}

func getValue(field reflect.StructField) string {
	var retval string
	defaultValue := field.Tag.Get("default")
	envField := field.Tag.Get("env")

	if os.Getenv(envField) != "" {
		retval = os.Getenv(envField)
	} else {
		retval = defaultValue
	}
	return retval
}

func setField(field reflect.Value, setVal string) error {
	if !field.CanSet() {
		return nil
	}

	if !shouldInitializeField(field, setVal) {
		return nil
	}

	isInitial := isInitialValue(field)
	if isInitial {
		switch field.Kind() {
		case reflect.Bool:
			if val, err := strconv.ParseBool(setVal); err == nil {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			}
		case reflect.Int:
			if val, err := strconv.ParseInt(setVal, 0, strconv.IntSize); err == nil {
				field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
			}
		case reflect.Int8:
			if val, err := strconv.ParseInt(setVal, 0, 8); err == nil {
				field.Set(reflect.ValueOf(int8(val)).Convert(field.Type()))
			}
		case reflect.Int16:
			if val, err := strconv.ParseInt(setVal, 0, 16); err == nil {
				field.Set(reflect.ValueOf(int16(val)).Convert(field.Type()))
			}
		case reflect.Int32:
			if val, err := strconv.ParseInt(setVal, 0, 32); err == nil {
				field.Set(reflect.ValueOf(int32(val)).Convert(field.Type()))
			}
		case reflect.Int64:
			if val, err := time.ParseDuration(setVal); err == nil {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			} else if val, err := strconv.ParseInt(setVal, 0, 64); err == nil {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			}
		case reflect.Uint:
			if val, err := strconv.ParseUint(setVal, 0, strconv.IntSize); err == nil {
				field.Set(reflect.ValueOf(uint(val)).Convert(field.Type()))
			}
		case reflect.Uint8:
			if val, err := strconv.ParseUint(setVal, 0, 8); err == nil {
				field.Set(reflect.ValueOf(uint8(val)).Convert(field.Type()))
			}
		case reflect.Uint16:
			if val, err := strconv.ParseUint(setVal, 0, 16); err == nil {
				field.Set(reflect.ValueOf(uint16(val)).Convert(field.Type()))
			}
		case reflect.Uint32:
			if val, err := strconv.ParseUint(setVal, 0, 32); err == nil {
				field.Set(reflect.ValueOf(uint32(val)).Convert(field.Type()))
			}
		case reflect.Uint64:
			if val, err := strconv.ParseUint(setVal, 0, 64); err == nil {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			}
		case reflect.Uintptr:
			if val, err := strconv.ParseUint(setVal, 0, strconv.IntSize); err == nil {
				field.Set(reflect.ValueOf(uintptr(val)).Convert(field.Type()))
			}
		case reflect.Float32:
			if val, err := strconv.ParseFloat(setVal, 32); err == nil {
				field.Set(reflect.ValueOf(float32(val)).Convert(field.Type()))
			}
		case reflect.Float64:
			if val, err := strconv.ParseFloat(setVal, 64); err == nil {
				field.Set(reflect.ValueOf(val).Convert(field.Type()))
			}
		case reflect.String:
			field.Set(reflect.ValueOf(setVal).Convert(field.Type()))

		case reflect.Slice:
			ref := reflect.New(field.Type())
			ref.Elem().Set(reflect.MakeSlice(field.Type(), 0, 0))
			if setVal != "" && setVal != "[]" {
				if err := json.Unmarshal([]byte(setVal), ref.Interface()); err != nil {
					return err
				}
			}
			field.Set(ref.Elem().Convert(field.Type()))
		case reflect.Map:
			ref := reflect.New(field.Type())
			ref.Elem().Set(reflect.MakeMap(field.Type()))
			if setVal != "" && setVal != "{}" {
				if err := json.Unmarshal([]byte(setVal), ref.Interface()); err != nil {
					return err
				}
			}
			field.Set(ref.Elem().Convert(field.Type()))
		case reflect.Struct:
			if setVal != "" && setVal != "{}" {
				if err := json.Unmarshal([]byte(setVal), field.Addr().Interface()); err != nil {
					return err
				}
			}
		case reflect.Ptr:
			field.Set(reflect.New(field.Type().Elem()))
		}
	}

	switch field.Kind() {
	case reflect.Ptr:
		if isInitial || field.Elem().Kind() == reflect.Struct {
			setField(field.Elem(), setVal)
		}
	case reflect.Struct:
		if err := SetDefaults(field.Addr().Interface()); err != nil {
			return err
		}
	case reflect.Slice:
		for j := 0; j < field.Len(); j++ {
			if err := setField(field.Index(j), setVal); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, e := range field.MapKeys() {
			var v = field.MapIndex(e)

			switch v.Kind() {
			case reflect.Ptr:
				switch v.Elem().Kind() {
				case reflect.Struct, reflect.Slice, reflect.Map:
					if err := setField(v.Elem(), ""); err != nil {
						return err
					}
				}
			case reflect.Struct, reflect.Slice, reflect.Map:
				ref := reflect.New(v.Type())
				ref.Elem().Set(v)
				if err := setField(ref.Elem(), ""); err != nil {
					return err
				}
				field.SetMapIndex(e, ref.Elem().Convert(v.Type()))
			}
		}
	}

	return nil
}

func PrettyPrint(v interface{}) {
	output, _ := json.MarshalIndent(v, "", "    ")
	fmt.Printf("%s\n", output)
}

func isInitialValue(field reflect.Value) bool {
	return reflect.DeepEqual(reflect.Zero(field.Type()).Interface(), field.Interface())
}

func shouldInitializeField(field reflect.Value, tag string) bool {
	switch field.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			return true
		}
	case reflect.Slice:
		return field.Len() > 0 || tag != ""
	case reflect.Map:
		return field.Len() > 0 || tag != ""
	}

	return tag != ""
}
