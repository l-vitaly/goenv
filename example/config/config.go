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
