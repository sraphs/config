package env_test

import (
	"fmt"
	"os"

	"github.com/sraphs/config"
	"github.com/sraphs/config/env"
	"github.com/sraphs/config/internal/testdata"
)

var conf = new(testdata.Conf)

func Example() {
	const prefix = "sraph_"
	envs := map[string]string{
		"sraph_log_level":               "info",
		"sraph_log_file":                "",
		"sraph_log_logger":              "stdlog",
		"sraph_server_http_addr":        ":8000",
		"sraph_server_http_timeout":     "1s",
		"sraph_server_grpc_addr":        ":9000",
		"sraph_server_grpc_timeout":     "1s",
		"sraph_data_database_driver":    "mysql",
		"sraph_data_database_dsn":       "root:root@tcp(mysql:3306)/test",
		"sraph_data_redis_addr":         "redis:6379",
		"sraph_data_redis_readTimeout":  "2s",
		"sraph_data_redis_writeTimeout": "3s",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}

	c := config.New(
		config.WithSource(env.NewSource(prefix)),
	)

	if err := c.Load(); err != nil {
		fmt.Println(err)
	}

	if err := c.Scan(&conf); err != nil {
		fmt.Println(err)
	}

	c.Watch(func(c config.Config) {
		fmt.Println(c)
	})

	fmt.Println(conf)

	// Outputs:
	// log:{level:"info"}  server:{http:{addr:":8000"  timeout:{seconds:1}}  grpc:{addr:":9000"  timeout:{seconds:1}}}  data:{database:{driver:"mysql"}  redis:{addr:"redis:6379"  read_timeout:{seconds:2}  write_timeout:{seconds:3}}}
}
