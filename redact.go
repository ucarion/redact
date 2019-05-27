// Package redact replaces elements of a struct with its zero value.
package redact

import (
	"fmt"
	"reflect"
)

// Redact replaces the (possibly nested) field of v at path with its zero value.
// Panics if v does not have the given path, or is otherwise unassignable.
// Panics if v is not a pointer.
//
// When Redact encounters an array or slice, every element of that array or
// slice will be recurisvely redacted. For example, if v is an array of users,
// then you can redact the passwords of all users in the array.
//
// When Redact encounters a struct or map, only the named element from path is
// redacted. Redact does not support updating all elements of a map, nor does it
// support updating the keys of a map.
//
// When Redact encounters a pointer, it recursively redacts the value that the
// pointer dereferences to.
//
// Go does not support mutating map elements. As a result, you have two options
// when dealing with maps:
//
// If you want to simply set the elements of a map to their zero values, Redact
// will do that for you. Just pass a path that points to an element of a map.
//
// If you want to mutate one of the elements of a map, then you have to map the
// values of the map be pointers. This is a limitation from Go itself, not this
// package.
func Redact(path []string, v interface{}) {
	redact(path, reflect.ValueOf(v).Elem())
}

func redact(path []string, v reflect.Value) {
	if len(path) == 0 {
		v.Set(reflect.Zero(v.Type()))
	} else {
		switch v.Type().Kind() {
		case reflect.Struct:
			redact(path[1:], v.FieldByName(path[0]))
		case reflect.Map:
			// You can't mutate elements of a map, but as a special case if we just
			// want to zero out one of the elements of the map, we can just do that
			// here.
			if len(path) == 1 {
				v.SetMapIndex(reflect.ValueOf(path[0]), reflect.Zero(v.Type().Elem()))
			} else {
				redact(path[1:], v.MapIndex(reflect.ValueOf(path[0])))
			}
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				redact(path, v.Index(i))
			}
		case reflect.Ptr:
			redact(path, v.Elem())
		default:
			panic(fmt.Sprintf("redact.Redact: unsupported type %v", v.Type().String()))
		}
	}
}
