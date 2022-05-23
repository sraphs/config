package flag

import (
	"fmt"
	"os"
	"testing"
)

func Test_flag(t *testing.T) {
	os.Args = []string{"", "--foo", "--bar=baz", "--sp", "2", "--a.b.c=d"}

	f := NewSource()

	ds, err := f.Load()

	if err != nil {
		t.Error(err)
	}

	var m = struct {
		Foo bool
		Bar string
		Sp  int
		A   struct {
			B struct {
				C string
			}
		}
	}{}

	for _, d := range ds {
		d.GetCodec().Unmarshal(d.Data, &m)
	}

	fmt.Println(m)

}
