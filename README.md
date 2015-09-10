# TypeRegistry

[![Build Status](https://travis-ci.org/rcarver/typeregistry.svg)](https://travis-ci.org/rcarver/typeregistry)
[![GoDoc](https://godoc.org/github.com/rcarver/typeregistry?status.svg)](https://godoc.org/github.com/rcarver/typeregistry)

TypeRegistry is a generic system to instantiate arbitrary types by name. Since
go cannot instantiate a type directly, we must first register any type that we
would later like to instantiate. The registry handles these mechanics for you.

In addition, the registry supports marshal, unmarshal and dependency injection
for getting objects in and out of storage.

## Example


To use, create a registry, then register add any types to it. Make sure to add
the pointer receiver if that's what you want to instantiate.

```golang
// The type you'd like to instantiate by name.
type simpleThing struct { }

// A registry.
registry := typeregistry.New()

// Register a type, this one by pointer receiver.
name := registry.Add(&simpleThing{})

// Get a *simpleThing
thing := registry.New(name)
```

See [this example file](example_unmarshal_test.go) for more detailed examples of marshal, unmarshal and custom setup.

## Common Usage

A common pattern for using this library is as a global handler to marshal/unmarshal a specific type. Here, implementations of a fictional `Conversation` type can be registered and then pulled in and out of storage formats. The package-level wrapper functions perform the typecasting needed to keep things simple for users.

```golang
package main

import "github.com/rcarver/typeregistry"

// Conversation is implemented many different ways.
type Conversation interface {
  Talk()
}

var registry = typeregistry.New()

// Register adds a new type to the conversation registry.
func Register(c Conversation) {
	registry.Add(c)
}

// Marshal encodes the conversation, returning its registered name, and the
// encoded bytes.
func Marshal(c Conversation) (string, []byte, error) {
	return registry.Marshal(c)
}

// ConvoDepsFunc is used to setup a conversation before it's unmarshaled.
type ConvoDepsFunc func(Conversation)

// Unmarshal decodes a conversation.
func Unmarshal(name string, data []byte, deps ConvoDepsFunc) (Conversation, error) {
	o, err := registry.Unmarshal(name, data, func(o interface{}) {
		if deps != nil {
			deps(o.(Conversation))
		}
	})
	if err == nil {
		return o.(Conversation), nil
	}
	return nil, err
}
```

## Author

Ryan Carver (ryan@ryancarver.com)

## License

MIT
