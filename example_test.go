package config_test

import (
	"path"

	"github.com/sraphs/go/config"
	"github.com/sraphs/go/config/env"
	"github.com/sraphs/go/config/file"
	"github.com/sraphs/go/config/flag"
	testData "github.com/sraphs/go/config/internal/testdata"
)

func Example() {
	p := path.Join("internal", "testdata")

	c := config.New(
		config.WithSource(
			env.NewSource("sraph_"),
			file.NewSource(p),
			flag.NewSource(),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var conf testData.Conf

	if err := c.Scan(&conf); err != nil {
		panic(err)
	}

	c.Watch(func(c config.Config) {
		c.Scan(&conf)
	})

	// fmt.Println(conf)

	// Output:
}
