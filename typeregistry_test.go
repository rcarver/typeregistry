package typeregistry

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type nothingType struct {
}

type nameType struct {
	Name string
}

type marshalBinType struct {
	Name string
	Fail bool
}

func (m marshalBinType) MarshalBinary() ([]byte, error) {
	if m.Fail {
		return nil, fmt.Errorf("Failed")
	}
	return []byte("bin:" + m.Name), nil
}

type unmarshalBinType struct {
	Name string
}

func (m *unmarshalBinType) UnmarshalBinary(data []byte) error {
	m.Name = "bin:" + string(data)
	return nil
}

type unmarshalBinFailType struct {
}

func (m *unmarshalBinFailType) UnmarshalBinary(data []byte) error {
	return fmt.Errorf("Failed")
}

type marshalTextType struct {
	Name string
	Fail bool
}

func (m marshalTextType) MarshalText() ([]byte, error) {
	if m.Fail {
		return nil, fmt.Errorf("Failed")
	}
	return []byte("text:" + m.Name), nil
}

type unmarshalTextType struct {
	Name string
}

func (m *unmarshalTextType) UnmarshalText(data []byte) error {
	m.Name = "bin:" + string(data)
	return nil
}

type unmarshalTextFailType struct {
}

func (m *unmarshalTextFailType) UnmarshalText(data []byte) error {
	return fmt.Errorf("Failed")
}

func TestNew(t *testing.T) {
	r := New()
	if len(r) != 0 {
		t.Errorf("New want empty, got %d", len(r))
	}
}

func TestTypeRegistry_Add(t *testing.T) {
	tests := []struct {
		t    interface{}
		want string
	}{
		{
			t:    nothingType{},
			want: "typeregistry.nothingType",
		},
		{
			t:    &nothingType{},
			want: "*typeregistry.nothingType",
		},
	}
	for i, test := range tests {
		r := make(TypeRegistry)
		got := r.Add(test.t)
		if got != test.want {
			t.Errorf("%d Add(%#v) got %s, want %s", i, test.t, got, test.want)
		}

	}
	var paniced string
	func() {
		r := make(TypeRegistry)
		defer func() {
			if r := recover(); r != nil {
				paniced = r.(string)
			}
		}()
		r.Add(nil)

	}()
	if paniced != "typeregistry cannot add nil" {
		t.Errorf("Expected Add(nil) to panic, got %s", paniced)
	}
}

func TestTypeRegistry_New(t *testing.T) {
	tests := []struct {
		t    interface{}
		want interface{}
	}{
		{
			t:    nothingType{},
			want: nothingType{},
		},
		{
			t:    &nothingType{},
			want: &nothingType{},
		},
		{
			t:    &nameType{"Hi"},
			want: &nameType{""},
		},
	}
	for i, test := range tests {
		r := make(TypeRegistry)
		name := r.Add(test.t)
		got := r.New(name)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d New(%s) got %#v, want %#v", i, name, got, test.want)
		}
	}
	var paniced string
	func() {
		r := make(TypeRegistry)
		defer func() {
			if r := recover(); r != nil {
				paniced = r.(string)
			}
		}()
		r.Add(nameType{})
		r.New("foo")
	}()
	if paniced != "typeregistry does not know \"foo\"" {
		t.Errorf("Expected New(\"foo\") to panic, got %s", paniced)
	}
}

func TestTypeRegistry_Marshal(t *testing.T) {
	tests := []struct {
		marsh interface{}
		name  string
		val   []byte
		err   bool
	}{
		{
			marsh: nothingType{},
			name:  "typeregistry.nothingType",
			val:   []byte{},
			err:   false,
		},
		{
			marsh: &nothingType{},
			name:  "*typeregistry.nothingType",
			val:   []byte{},
			err:   false,
		},
		{
			marsh: marshalBinType{Name: "ok", Fail: false},
			name:  "typeregistry.marshalBinType",
			val:   []byte("bin:ok"),
			err:   false,
		},
		{
			marsh: marshalBinType{Name: "ok", Fail: true},
			name:  "typeregistry.marshalBinType",
			val:   []byte{},
			err:   true,
		},
		{
			marsh: marshalTextType{Name: "ok", Fail: false},
			name:  "typeregistry.marshalTextType",
			val:   []byte("text:ok"),
			err:   false,
		},
		{
			marsh: marshalTextType{Name: "ok", Fail: true},
			name:  "typeregistry.marshalTextType",
			val:   []byte{},
			err:   true,
		},
	}
	for i, test := range tests {
		r := make(TypeRegistry)
		name, val, err := r.Marshal(test.marsh)
		if name != test.name {
			t.Errorf("%d Marshal() name got %#v, want %#v", i, name, test.name)
		}
		if test.err {
			if err == nil {
				t.Errorf("%d Marshal() wants error, got none", i)
			}
		} else {
			if err != nil {
				t.Errorf("%d Marshal() wants no error, got: %s", i, err)
			}
		}
		if !bytes.Equal(val, test.val) {
			t.Errorf("%d Marshal() value: got %#v, want %#v", i, val, test.val)
		}
	}
}

func TestTypeRegistry_Unmarshal(t *testing.T) {
	tests := []struct {
		t    interface{}
		data []byte
		deps DepsFunc
		err  bool
		want interface{}
	}{
		{
			t:    nothingType{},
			data: []byte{},
			deps: NoDeps,
			err:  false,
			want: nothingType{},
		},
		{
			t:    &nothingType{},
			data: []byte{},
			deps: NoDeps,
			err:  false,
			want: &nothingType{},
		},
		{
			t:    &unmarshalBinType{},
			data: []byte("ok"),
			deps: NoDeps,
			err:  false,
			want: &unmarshalBinType{Name: "bin:ok"},
		},
		{
			t:    &unmarshalBinFailType{},
			data: []byte("ok"),
			deps: NoDeps,
			err:  true,
			want: &unmarshalBinFailType{},
		},
		{
			t:    &unmarshalTextType{},
			data: []byte("ok"),
			deps: NoDeps,
			err:  false,
			want: &unmarshalTextType{Name: "bin:ok"},
		},
		{
			t:    &unmarshalTextFailType{},
			data: []byte("ok"),
			deps: NoDeps,
			err:  true,
			want: &unmarshalTextFailType{},
		},
		{
			t:    &nameType{},
			data: []byte{},
			deps: func(i interface{}) {
				if x, ok := i.(*nameType); ok {
					x.Name = "ok"
				}
			},
			err:  false,
			want: &nameType{"ok"},
		},
	}
	for i, test := range tests {
		r := make(TypeRegistry)
		name := r.Add(test.t)
		got, err := r.Unmarshal(name, test.data, test.deps)
		if test.err {
			if err == nil {
				t.Errorf("%d Unmarshal wants error, got none", i)
			}
		} else {
			if err != nil {
				t.Errorf("%d Unmarshal wants no error, got: %s", i, err)
			}
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d Unmarshal() got %#v, want %#v", i, got, test.want)
		}
	}
}
