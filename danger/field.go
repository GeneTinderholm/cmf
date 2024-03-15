package danger

import "reflect"

// SetPrivateField is dangerous (hence the package name)
// don't use it... but, like, sometimes there's no other option
func SetPrivateField[S, T any](thing *S, fieldName string, val T) {
	ptr := reflect.ValueOf(thing).Elem().FieldByName(fieldName).Addr().UnsafePointer()
	*(*T)(ptr) = val
}
