package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/l-vitaly/goenv"
)

const envCfgTemplate = `
package config

import (
	"fmt"

	"github.com/l-vitaly/goenv"
)

// env name constants
const (
	{{range $key, $val := .}}
	{{$key}}EnvName = "{{$key}}"
	{{end}}
)

// Config service configuration
type Config struct {
}

// Get get env config vars
func Get() (*Config, error) {
	cfg := &Config{}
	goenv.StringVar(&cfg.HTTPAddr, HTTPAddrEnvName, ":9000")

	goenv.Parse()
	if cfg.Mongo.URL == "" {
		return nil, fmt.Errorf("could not set %s", DBConnStrEnvName)
	}
	return cfg, nil
}
`

func usage() string {
	return fmt.Sprintf("Usage: %s <filename> (try -h)", os.Args[0])
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal(usage())
	}
	filename := flag.Arg(0)
	l := new(goenv.Loader)

	envVars, err := l.Load(filename)
	if err != nil {
		log.Fatal(err)
	}

	t := template.Must(template.New("envCfg").Parse(envCfgTemplate))

	t.Execute(os.Stdout, envVars)
}
