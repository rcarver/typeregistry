package typeregistry_test

import (
	"encoding/json"
	"fmt"

	"github.com/rcarver/typeregistry"
)

// User is your typical user struct.
type User struct {
	ID   string
	Name string
}

// UserService is your typical service backend to retrieve users.
type UserService struct {
	users []*User
}

// FindUser returns a user by id.
func (s *UserService) FindUser(id string) *User {
	for _, u := range s.users {
		if u.ID == id {
			return u
		}
	}
	return nil
}

// userThing is a sample struct that can marshal/unmarshal itself and requires
// a service collaborator when doing so.
type userThing struct {
	UserID string
	user   *User
	svc    *UserService
}

func (t *userThing) Marshal() ([]byte, error) {
	// We store the ID when marshaling so we can look it up fresh when
	// unmarshaled.
	t.UserID = t.user.ID
	return json.Marshal(t)
}

func (t *userThing) Unmarshal(data []byte) error {
	// Unmarshal, getting the UserID.
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	// Use the UserID and injected service to restore the user.
	t.user = t.svc.FindUser(t.UserID)
	return nil
}

func ExampleTypeRegistry_Unmarshal() {
	registry := typeregistry.New()
	registry.Add(&userThing{})

	// The service world.
	var (
		ryan   = &User{"1", "Ryan"}
		svc    = &UserService{[]*User{ryan}}
		sample = &userThing{user: ryan}
	)

	// Pass between marshal and unmarshal.
	var (
		name string
		data []byte
	)

	// Marshal it.
	func() {
		var err error
		name, data, err = registry.Marshal(sample)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sample marshaled name:%s data:%s\n", name, data)
	}()

	// Unmarshal it with custom setup.
	func() {
		si, err := registry.Unmarshal(name, data, func(o interface{}) {
			if s, ok := o.(*userThing); ok {
				s.svc = svc
			}
		})
		if err != nil {
			panic(err)
		}
		sample := si.(*userThing)
		fmt.Printf("Sample unmarshaled with ID:%s, Name:%s\n", sample.UserID, sample.user.Name)
	}()
	// Output:
	// Sample marshaled name:*typeregistry_test.userThing data:{"UserID":"1"}
	// Sample unmarshaled with ID:1, Name:Ryan
}
