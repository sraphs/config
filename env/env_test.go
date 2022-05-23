package env

import (
	"os"
	"reflect"
	"testing"

	"github.com/sraphs/maps"

	"github.com/sraphs/config"
)

func TestEnvWithPrefix(t *testing.T) {
	// set env
	prefix := "sraph_"
	envs := map[string]string{
		prefix + "service_name": "sraph_app",
		prefix + "addr":         "192.168.0.1",
		prefix + "age":          "20",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}

	c := config.New(
		config.WithSource(NewSource(prefix)),
	)

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	conf := make(map[string]interface{})

	if err := c.Scan(&conf); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		path   string
		expect interface{}
	}{
		{
			name:   "service.name",
			path:   "service.name",
			expect: "sraph_app",
		},
		{
			name:   "addr",
			path:   "addr",
			expect: "192.168.0.1",
		},
		{
			name:   "age",
			path:   "age",
			expect: "20",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := maps.Get(conf, test.path).String()
			if !reflect.DeepEqual(test.expect, actual) {
				t.Errorf("expect %v, actual %v", test.expect, actual)
			}
		})
	}
}

func TestEnvWithoutPrefix(t *testing.T) {
	envs := map[string]string{
		"service_name": "sraph_app",
		"addr":         "192.168.0.1",
		"age":          "20",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}

	c := config.New(
		config.WithSource(NewSource("")),
	)

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	conf := make(map[string]interface{})

	if err := c.Scan(&conf); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		path   string
		expect interface{}
	}{
		{
			name:   "service.name",
			path:   "service.name",
			expect: "sraph_app",
		},
		{
			name:   "addr",
			path:   "addr",
			expect: "192.168.0.1",
		},
		{
			name:   "age",
			path:   "age",
			expect: "20",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := maps.Get(conf, test.path).String()
			if !reflect.DeepEqual(test.expect, actual) {
				t.Errorf("expect %v, actual %v", test.expect, actual)
			}
		})
	}
}
