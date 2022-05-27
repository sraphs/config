package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sraphs/encoding"
)

func TestReader_Merge(t *testing.T) {
	var (
		err error
		ok  bool
	)
	opts := options{
		decoder: func(d *Descriptor, v map[string]interface{}) error {
			if codec := encoding.GetCodec(d.Format); codec != nil {
				return codec.Unmarshal(d.Data, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", d.Name, d.Format)
		},
		resolver: defaultResolver,
	}
	r := newReader(opts)
	err = r.Merge(&Descriptor{
		Name:   "a",
		Data:   []byte("bad"),
		Format: "json",
	})
	if err == nil {
		t.Fatal(`err is nil`)
	}

	err = r.Merge(&Descriptor{
		Name:   "b",
		Data:   []byte(`{"nice": "boat", "x": 1}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(`err is not nil`, err)
	}
	vv, ok := r.Value("nice")
	if !ok {
		t.Fatal(`ok is false`)
	}
	vvv, err := vv.String()
	if err != nil {
		t.Fatal(`err is not nil`, err)
	}
	if vvv != "boat" {
		t.Fatal(`vvv is not equal to "boat"`)
	}

	err = r.Merge(&Descriptor{
		Name:   "b",
		Data:   []byte(`{"x": 2}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(`err is not nil`, err)
	}
	vv, ok = r.Value("x")
	if !ok {
		t.Fatal(`ok is false`)
	}
	vvx, err := vv.Int()
	if err != nil {
		t.Fatal(`err is not nil`, err)
	}
	if int64(2) != vvx {
		t.Fatal(`vvx is not equal to 2`)
	}
}

func TestReader_Value(t *testing.T) {
	opts := options{
		decoder: func(d *Descriptor, v map[string]interface{}) error {
			if codec := encoding.GetCodec(d.Format); codec != nil {
				return codec.Unmarshal(d.Data, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", d.Name, d.Format)
		},
		resolver: defaultResolver,
	}

	ymlval := `
a: 
  b: 
    X: 1
    Y: "lol"
    z: true
`
	tests := []struct {
		name string
		kv   Descriptor
	}{
		{
			name: "json value",
			kv: Descriptor{
				Name:   "config",
				Data:   []byte(`{"a": {"b": {"X": 1, "Y": "lol", "z": true}}}`),
				Format: "json",
			},
		},
		{
			name: "yaml value",
			kv: Descriptor{
				Name:   "config",
				Data:   []byte(ymlval),
				Format: "yaml",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := newReader(opts)
			err := r.Merge(&test.kv)
			if err != nil {
				t.Fatal(`err is not nil`, err)
			}
			vv, ok := r.Value("a.b.X")
			if !ok {
				t.Fatal(`ok is false`)
			}
			vvv, err := vv.Int()
			if err != nil {
				t.Fatal(`err is not nil`, err)
			}
			if int64(1) != vvv {
				t.Fatal(`vvv is not equal to 1`)
			}

			if err != nil {
				t.Fatal(`err is not nil`, err)
			}
			vv, ok = r.Value("a.b.Y")
			if !ok {
				t.Fatal(`ok is false`)
			}
			vvy, err := vv.String()
			if err != nil {
				t.Fatal(`err is not nil`, err)
			}
			if vvy != "lol" {
				t.Fatal(`vvy is not equal to "lol"`)
			}

			if err != nil {
				t.Fatal(`err is not nil`, err)
			}
			vv, ok = r.Value("a.b.z")
			if !ok {
				t.Fatal(`ok is false`)
			}
			vvz, err := vv.Bool()
			if err != nil {
				t.Fatal(`err is not nil`, err)
			}
			if !vvz {
				t.Fatal(`vvz is not equal to true`)
			}

			_, ok = r.Value("aasasdg=234l.asdfk,")
			if ok {
				t.Fatal(`ok is true`)
			}

			_, ok = r.Value("aas......asdg=234l.asdfk,")
			if ok {
				t.Fatal(`ok is true`)
			}

			_, ok = r.Value("a.b.Y.")
			if ok {
				t.Fatal(`ok is true`)
			}
		})
	}
}

func TestReader_Source(t *testing.T) {
	var err error
	opts := options{
		decoder: func(d *Descriptor, v map[string]interface{}) error {
			if codec := encoding.GetCodec(d.Format); codec != nil {
				return codec.Unmarshal(d.Data, &v)
			}
			return fmt.Errorf("unsupported key: %s format: %s", d.Name, d.Format)
		},
		resolver: defaultResolver,
	}
	r := newReader(opts)
	err = r.Merge(&Descriptor{
		Name:   "b",
		Data:   []byte(`{"a": {"b": {"X": 1}}}`),
		Format: "json",
	})
	if err != nil {
		t.Fatal(`err is not nil`, err)
	}
	b, err := r.Source()
	if err != nil {
		t.Fatal(`err is not nil`, err)
	}
	if !reflect.DeepEqual([]byte(`{"a":{"b":{"X":1}}}`), b) {
		t.Fatal("[]byte(`{\"a\":{\"b\":{\"X\":1}}}`) is not equal to b")
	}
}
