// Package typeregistry is a generic system to instantiate types by name. Since
// go cannot instantiate a type directly, we must first register any type that
// we would later like to instantiate. The registry handles these mechanics for
// you.
//
// In addition, the registry supports marshal, unmarshal, and custom setup
// injection for getting objects in and out of storage.
//
// Marshaling an object results in the registered name of the type, plus byte
// data if the type implements Marshaler.
//
// Unmarshaling performs the reverse operation, first instantiating the type by
// name, then using Unmarshaler (if implemented) to populate the object (note
// that the type should probably be a pointer reciever for this to be useful).
// If the object requires collaborators, or data from the outside world then a
// function can be passed to Unmarshal that receives the object after it's
// instantiated and before it's unmarshaled.
package typeregistry

import (
	"fmt"
	"reflect"
)

// Marshaler is implemented by any type that can encode a copy of itself. The
// style of encoding doesn't matter, it will only be seen by Unmarshaler.
type Marshaler interface {
	Marshal() ([]byte, error)
}

// Unmarshaler is implemented by any type that can decode of a copy of itself,
// as returned by its Marshal method.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

// TypeRegistry can instantiate, marshal, and unmarshal types from string names
// and type-defined encodings.
type TypeRegistry map[string]reflect.Type

// New initializes an empty TypeRegistry.
func New() TypeRegistry {
	return make(TypeRegistry)
}

// Add puts a new type in the registry. If the type cannot be registered, it
// panics. It returns the name that it was registered as.
func (r TypeRegistry) Add(o interface{}) string {
	if o == nil {
		panic("typeregistry cannot add nil")
	}
	name := r.name(o)
	r[name] = reflect.TypeOf(o)
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

// Marshal encodes a type. If the type implements Marshaler or its bytes are
// returned.
func (r TypeRegistry) Marshal(o interface{}) (string, []byte, error) {
	var (
		name  = r.name(o)
		bytes []byte
		err   error
	)
	switch m := o.(type) {
	case Marshaler:
		bytes, err = m.Marshal()
	}
	return name, bytes, err
}

// SetupFunc is passed to Unmarshal to manually manipulate the object after
// it's instantiated, but before it's unmarshaled. This can be used to set
// dependencies that are needed during unmarshal.  For example, to covert a
// user's ID to a user object.
type SetupFunc func(interface{})

// NoSetup is a SetupFunc that does nothing. It is functionally equivalent to
// passing nil, but it's more descriptive so please do.
var NoSetup = func(i interface{}) {}

// Unmarshal decodes a type by name. If the type implements Unmarshaler, the
// data is used to unmarshal. SetupFunc can be passed to inject any other data
// into the type before it is unmarshaled.
func (r TypeRegistry) Unmarshal(name string, data []byte, setup SetupFunc) (interface{}, error) {
	instance := r.New(name)
	if setup != nil {
		setup(instance)
	}
	switch m := instance.(type) {
	case Unmarshaler:
		if err := m.Unmarshal(data); err != nil {
			return instance, err
		}
	}
	return instance, nil
}

func (r TypeRegistry) name(c interface{}) string {
	// TODO: let types set their own name?
	return reflect.TypeOf(c).String()
}
