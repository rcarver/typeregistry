package typeregistry_test

import (
	"fmt"

	"github.com/rcarver/typegregistry"
)

type simpleThing struct {
	Name string
}

func ExampleTypeRegistry_New() {
	registry := typeregistry.New()
	name := registry.Add(simpleThing{})
	thing := registry.New(name)
	fmt.Printf("%#v", thing)
	// Output: typeregistry_test.simpleThing{Name:""}
}
