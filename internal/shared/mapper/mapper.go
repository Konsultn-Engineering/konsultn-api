package mapper

import (
	"github.com/jinzhu/copier"
	"reflect"
)

// Convert maps from type U to T, then removes any zero-valued nested pointers.
func Convert[T any, U any](from U) (T, error) {
	var to T
	if err := copier.Copy(&to, &from); err != nil {
		return to, err
	}
	cleanEmptyPointers(&to)
	return to, nil
}

func cleanEmptyPointers(v any) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}
	cleanStruct(val.Elem())
}

func cleanStruct(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if !f.CanSet() {
			continue
		}

		switch f.Kind() {
		case reflect.Ptr:
			if f.Type().Elem().Kind() == reflect.Struct && !f.IsNil() {
				cleanStruct(f.Elem())
				if f.Elem().IsZero() {
					f.Set(reflect.Zero(f.Type()))
				}
			}
		case reflect.Struct:
			cleanStruct(f)

		case reflect.Slice:
			for j := 0; j < f.Len(); j++ {
				elem := f.Index(j)
				if elem.Kind() == reflect.Ptr && elem.Elem().Kind() == reflect.Struct {
					cleanStruct(elem.Elem())
				} else if elem.Kind() == reflect.Struct {
					cleanStruct(elem)
				}
			}
		default:
			break
		}
	}
}
