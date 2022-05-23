# config

[![CI](https://github.com/sraphs/config/actions/workflows/ci.yml/badge.svg)](https://github.com/sraphs/config/actions/workflows/ci.yml)

>  Go configuration library

## Usage

```go
package config_test

import (
	"path"

	"github.com/sraphs/config"
	"github.com/sraphs/config/env"
	"github.com/sraphs/config/file"
	"github.com/sraphs/config/flag"
	testData "github.com/sraphs/config/internal/testdata"
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

```

## Contributing

We alway welcome your contributions :clap:

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request


## CHANGELOG
See [Releases](https://github.com/sraphs/config/releases)

## License
[MIT © sraph.com](./LICENSE)
