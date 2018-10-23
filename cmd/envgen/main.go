package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/l-vitaly/goenv/parser"
)

const envCfgTemplate = `
{{define "const"}}{{if not .Child}}{{.Const}} = "{{.Env}}" ` + "\n" + ` {{else}}{{template "const" .Child}}{{end}}{{end}}
{{define "struct"}}{{if not .Child}}{{.Name}} {{.Type}}` + "\n" + `{{else}}{{.Name}} struct {` + "\n" + `{{template "struct" .Child}}` + "\n" + `}` + "\n" + `{{end}}{{end}}
{{define "func"}}{{if not .Child}}goenv.{{.Func}}(&cfg{{.CfgPath}}.{{.Name}}, {{.Const}}, {{.Value}})` + "\n" + `{{else}}{{template "func" .Child}}{{end}}{{end}}
{{define "valid"}}{{if not .Child}}if cfg{{.CfgPath}}.{{.Name}} == {{.EmptyValue}} {` + "\n" + `return nil, fmt.Errorf(errPattern, {{.Const}})}` + "\n" + `{{else}}{{template "valid" .Child}}{{end}}{{end}}

package config

import (
	"fmt"

	"github.com/l-vitaly/goenv"
)

var urlNil = url.URL{}

const errPattern = "could not set %s"

// env name constants
const (
{{range $node := .}}{{template "const" $node}}{{end}}
)

// Config config.
type Config struct {
{{range $node := .}}{{template "struct" $node}}{{end}}
}

// Parse env config vars.
func Parse() (*Config, error) {
	cfg := &Config{}

	{{range $node := .}}{{template "func" $node}}{{end}}
	goenv.Parse()

	{{range $node := .}}{{template "valid" $node}}{{end}}

	return cfg, nil
}
`

type node struct {
	Name       string
	Type       string
	Value      interface{}
	RawValue   interface{}
	EmptyValue interface{}
	Func       string
	Env        string
	Const      string
	CfgPath    string
	Child      *node
}

func usage() string {
	return fmt.Sprintf("Usage: %s <filename> <prefix> (try -h)", os.Args[0])
}

func normalize(s string, firstUpper bool) string {
	buf := new(bytes.Buffer)

	isUpper := firstUpper

	for _, r := range s {
		if r == '_' {
			isUpper = true
			continue
		}
		if !isUpper && r >= 'A' && r <= 'Z' {
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
			isUpper = false
		}
	}
	return buf.String()
}

func normalizeValue(t interface{}) interface{} {
	switch v := t.(type) {
	case *url.URL:
		return "url.URL{}"
	case time.Duration:
		return strconv.Itoa(int(v.Seconds())) + " * time.Second"
	case bool:
		return strconv.FormatBool(v)
	case uint64:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.Itoa(int(v))
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case string:
		return template.HTML("\"" + v + "\"")
	default:
		return template.HTML("\"\"")
	}
}

func valType(t interface{}) string {
	switch t.(type) {
	case bool:
		return "bool"
	case uint64:
		return "uint"
	case int64:
		return "int"
	case float64:
		return "float"
	case *url.URL:
		return "url.URL"
	case time.Duration:
		return "time.Duration"
	default:
		return "string"
	}
}

func emptyVal(t interface{}) interface{} {
	switch t.(type) {
	case *url.URL:
		return "urlNil"
	case uint64, int64, float64, time.Duration:
		return "0"
	default:
		return template.HTML("\"\"")
	}
}

func envFunc(t interface{}) string {
	switch t.(type) {
	case *url.URL:
		return "URLVar"
	case time.Duration:
		return "DurationVar"
	case bool:
		return "BoolVar"
	case uint64:
		return "UintVar"
	case float64:
		return "FloatVar"
	case int64:
		return "IntVar"
	default:
		return "StringVar"
	}
}

func traverse(parents []string, value *node) *node {
	if len(parents) == 0 {
		return value
	}
	root := &node{Name: normalize(parents[0], true)}
	next := root
	cfgPath := root.Name

	for i := 1; i < len(parents); i++ {
		next.Child = &node{Name: normalize(parents[i], true)}
		cfgPath = cfgPath + "." + next.Child.Name
		next = next.Child
	}

	value.CfgPath = "." + cfgPath
	next.Child = value
	return root
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal(usage())
	}
	filename := flag.Arg(0)

	prefix := ""
	if len(os.Args) == 3 {
		prefix = flag.Arg(1)
	}

	l := parser.New()

	envVars, err := l.Parse(filename)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	nodes := []*node{}

	for _, envVar := range envVars {
		envName := envVar.Name
		name := strings.TrimPrefix(envName, prefix)
		parents := strings.Split(name, "__")
		normalizeName := normalize(parents[len(parents)-1], true)
		constName := normalizeName + "EnvName"

		v := &node{
			Name:       normalizeName,
			Const:      constName,
			Env:        envName,
			Type:       valType(envVar.Value),
			Value:      normalizeValue(envVar.Value),
			RawValue:   envVar.Value,
			Func:       envFunc(envVar.Value),
			EmptyValue: emptyVal(envVar.Value),
		}

		nodes = append(nodes, traverse(parents[:len(parents)-1], v))
	}

	for _, n := range nodes {
		if n.Child != nil {
			fmt.Println(n.Child.Child)
		}
	}

	t := template.Must(template.New("envCfg").Parse(envCfgTemplate))

	buf := new(bytes.Buffer)

	t.Execute(buf, nodes)

	result, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(err, buf.String())
		os.Exit(1)
	}

	fmt.Println(string(result))
}
