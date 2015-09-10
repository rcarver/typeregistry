// Package typeregistry is a generic system to instantiate arbitrary types by
// name. Since go cannot instantiate a type directly, we must first register
// any type that we would later like to instantiate. The registry handles these
// mechanics for you.
//
// In addition, the registry supports marshal, unmarshal and dependency
// injection patterns.

// Marshaling an object results in the registered name of the type, plus byte
// data if the type implements encoding.BinaryMarshaler or
// encoding.TextMarshaler.
//
// Unmarshal performs the reverse operation, first instantiating the type by
// name, then using encoding.BinaryUnmarshaler or encoding.TextUnmarshaler to
// populate the object (note that the type should probably be a pointer
// reciever for this to be useful). If the object requires collaborators, or
// data from the outside world then a function can be passed to Unmarshal that
// receives the object after it's instantiated and before it's unmarshaled.
package typeregistry

import (
	"encoding"
	"fmt"
	"reflect"
)

// TypeRegistry can instantiate, marshal, and unmarshal types from string names
// and arbitrary encodings.
type TypeRegistry map[string]reflect.Type

// New initializes an empty TypeRegistry.
func New() TypeRegistry {
	return make(TypeRegistry)
}

// Add includes a new type in the registry. If the type cannot be registered,
// it panics. It returns the name that was registered.
func (r TypeRegistry) Add(c interface{}) string {
	if c == nil {
		panic("typeregistry cannot add nil")
	}
	name := r.name(c)
	r[name] = reflect.TypeOf(c)
	return name
}

// New instantiates a type by name. If the name is unknown, it panics.
func (r TypeRegistry) New(name string) interface{} {
	if val, ok := r[name]; ok {
		if val.Kind() == reflect.Ptr {
			v := reflect.New(val.Elem())
			return v.Interface()
		}
		v := reflect.New(val).Elem()
		return v.Interface()
	}
	panic(fmt.Sprintf("typeregistry does not know %#v", name))
}

// Marshal encodes a type. If the type implements encoding.BinaryMarshaler or
// encoding.TextMarshaler, its bytes are returned.
func (r TypeRegistry) Marshal(c interface{}) (string, []byte, error) {
	var (
		name  = r.name(c)
		bytes []byte
		err   error
	)
	switch m := c.(type) {
	case encoding.BinaryMarshaler:
		bytes, err = m.MarshalBinary()
	case encoding.TextMarshaler:
		bytes, err = m.MarshalText()
	}
	return name, bytes, err
}

// DepsFunc is passed to Unmarshal to manually manipulate the object after it's
// instantiated, but before it's unmarshaled. This can be used to set
// dependencies that are needed during unmarshal.  For example, to covert a
// user's ID to a user object.
type DepsFunc func(interface{})

// NoDeps is a DepsFunc that does nothing. It is functionally equivalent to
// passing nil, but it's more descriptive so please do.
var NoDeps = func(i interface{}) {}

// Unmarshal decodes a type by name. If the type implements
// encoding.BinaryUnmarshaler or encoding.TextUnmarshaler, the data is used to
// unmarshal. DepsFunc can be passed to inject any other data into the type
// before it is unmarshaled.
func (r TypeRegistry) Unmarshal(name string, data []byte, deps DepsFunc) (interface{}, error) {
	instance := r.New(name)
	if deps != nil {
		deps(instance)
	}
	switch m := instance.(type) {
	case encoding.BinaryUnmarshaler:
		if err := m.UnmarshalBinary(data); err != nil {
			return instance, err
		}
	case encoding.TextUnmarshaler:
		if err := m.UnmarshalText(data); err != nil {
			return instance, err
		}
	}
	return instance, nil
}

func (r TypeRegistry) name(c interface{}) string {
	// TODO: let types set their own name?
	return reflect.TypeOf(c).String()
}
