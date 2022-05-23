package file

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/sraphs/maps"

	"github.com/sraphs/config"
)

const (
	_testJSON = `
{
    "test":{
        "settings":{
            "int_key":1000,
            "float_key":1000.1,
            "duration_key":10000,
            "string_key":"string_value"
        },
        "server":{
            "addr":"127.0.0.1",
            "port":8000
        }
    },
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ]
}`

	_testJSONUpdate = `
{
    "test":{
        "settings":{
            "int_key":1000,
            "float_key":1000.1,
            "duration_key":10000,
            "string_key":"string_value"
        },
        "server":{
            "addr":"127.0.0.1",
            "port":8000
        }
    },
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ],
	"bar":{
		"event":"update"
	}
}`
)

type testSettings struct {
	IntKey      int64         `json:"int_key"`
	FloatKey    float64       `json:"float_key"`
	StringKey   string        `json:"string_key"`
	DurationKey time.Duration `json:"duration_key"`
}

var settings testSettings

func TestFile(t *testing.T) {
	var (
		path = filepath.Join(t.TempDir(), "test_config")
		file = filepath.Join(path, "test.json")
		data = []byte(_testJSON)
	)

	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(file, data, 0o666); err != nil {
		t.Error(err)
	}

	testSource(t, file, data)
	testSource(t, path, data)

	testWatchFile(t, file)
	testWatchDir(t, path, file)
}

func testSource(t *testing.T, path string, data []byte) {
	t.Log(path)
	s := NewSource(path)
	ds, err := s.Load()
	if err != nil {
		t.Error(err)
	}

	if string(ds[0].Data) != string(data) {
		t.Errorf("no expected: %s, but got: %s", ds[0].Data, data)
	}
}

func testWatchFile(t *testing.T, path string) {
	t.Log(path)

	s := NewSource(path)
	watch, err := s.Watch()
	if err != nil {
		t.Error(err)
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	_, err = f.WriteString(_testJSONUpdate)
	if err != nil {
		t.Error(err)
	}

	ds, err := watch.Next()
	if err != nil {
		t.Errorf(`watch.Next() error(%v)`, err)
	}
	if !reflect.DeepEqual(string(ds[0].Data), _testJSONUpdate) {
		t.Errorf(`string(ds[0].Data(%v) is  not equal to _testJSONUpdate(%v)`, ds[0].Data, _testJSONUpdate)
	}

	newFilepath := filepath.Join(filepath.Dir(path), "test1.json")
	if err = os.Rename(path, newFilepath); err != nil {
		t.Error(err)
	}
	ds, err = watch.Next()
	if err == nil {
		t.Errorf(`watch.Next() error(%v)`, err)
	}
	if ds != nil {
		t.Errorf(`watch.Next() error(%v)`, err)
	}

	err = watch.Stop()
	if err != nil {
		t.Errorf(`watch.Stop() error(%v)`, err)
	}

	if err := os.Rename(newFilepath, path); err != nil {
		t.Error(err)
	}
}

func testWatchDir(t *testing.T, path, file string) {
	t.Log(path)
	t.Log(file)

	s := NewSource(path)
	watch, err := s.Watch()
	if err != nil {
		t.Error(err)
	}

	f, err := os.OpenFile(file, os.O_RDWR, 0)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	_, err = f.WriteString(_testJSONUpdate)
	if err != nil {
		t.Error(err)
	}

	ds, err := watch.Next()
	if err != nil {
		t.Errorf(`watch.Next() error(%v)`, err)
	}
	if !reflect.DeepEqual(string(ds[0].Data), _testJSONUpdate) {
		t.Errorf(`string(ds[0].Data(%v) is  not equal to _testJSONUpdate(%v)`, ds[0].Data, _testJSONUpdate)
	}
}

func TestConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test_config.json")
	defer os.Remove(path)
	if err := os.WriteFile(path, []byte(_testJSON), 0o666); err != nil {
		t.Error(err)
	}

	c := config.New(
		config.WithSource(NewSource(path)),
	)

	testConfig(t, c)
}

func testConfig(t *testing.T, c config.Config) {
	expected := map[string]interface{}{
		"test.settings.int_key":      int64(1000),
		"test.settings.float_key":    float64(1000.1),
		"test.settings.string_key":   "string_value",
		"test.settings.duration_key": time.Duration(10000),
		"test.server.addr":           "127.0.0.1",
		"test.server.port":           int64(8000),
	}

	if err := c.Load(); err != nil {
		t.Error(err)
	}

	conf := make(map[string]interface{})

	// scan
	if err := c.Scan(&conf); err != nil {
		t.Error(err)
	}

	for key, value := range expected {
		switch value.(type) {
		case int64:
			if v, err := maps.Get(conf, key).Int(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case float64:
			if v, err := maps.Get(conf, key).Float(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case string:
			if v, err := maps.Get(conf, key).String(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case time.Duration:
			if v, err := maps.Get(conf, key).Duration(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		}
	}
}
