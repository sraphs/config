package config_test

import (
	"fmt"
	"os"
	"path"

	"github.com/sraphs/config"
	"github.com/sraphs/config/env"
	"github.com/sraphs/config/file"
	"github.com/sraphs/config/flag"
	testData "github.com/sraphs/config/internal/testdata"
)

func Example() {
	os.Setenv("sraph_log_level", "warn")

	p := path.Join("internal", "testdata", "yaml")

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

	fmt.Println(&conf)

	c.Get("log.level").String()

	// Outputs:
	// log:{level:"warn"}  server:{http:{addr:"0.0.0.0:8000"  timeout:{seconds:1}}  grpc:{addr:"0.0.0.0:9000"  timeout:{seconds:1}}}  data:{database:{driver:"mysql"}  redis:{addr:"mysql:6379"  read_timeout:{nanos:200000000}  write_timeout:{nanos:200000000}}}
}
