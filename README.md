# TypeRegistry

TypeRegistry is a generic system to instantiate arbitrary types by name. Since
go cannot instantiate a type directly, we must first register any type that we
would later like to instantiate. The registry handles these mechanics for you.

In addition, the registry supports marshal, unmarshal and dependency injection
for getting objects in and out of storage.

## Example


To use, create a registry, then register add any types to it. Make sure to add
the pointer receiver if that's what you want to instantiate.

```golang
# The type you'd like to instantiate by name.
type simpleThing struct { }

# A registry.
registry := typeregistry.New()

// Register a type, this one by pointer receiver.
name := registry.Add(&simpleThing{})

# Get a *simpleThing
thing := registry.New(name)
```

See the example files for more details.

## Author

Ryan Carver (ryan@ryancarver.com)

## License

MIT
