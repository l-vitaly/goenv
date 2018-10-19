# Example

main.go

```golang
package main

import (
	"log"

	"github.com/l-vitaly/goenv/example/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg.DBConn)
}
```

config/config.go

```golang
package config

import (
	"fmt"

	"github.com/l-vitaly/goenv"
)

const errPattern = "could not set %s"

// env name consts
const (
	DBConnEnvName = "DB_CONN"
)

// Config config.
type Config struct {
	DBConn string
}

// Parse env config vars.
func Parse() (*Config, error) {
	cfg := &Config{}

	goenv.StringVar(&cfg.DBConn, DBConnEnvName, "")
	goenv.Parse()

	if cfg.DBConn == "" {
		return nil, fmt.Errorf(errPattern, DBConnEnvName)
	}

	return cfg, nil
}
```
