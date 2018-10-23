package parser

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	lineRegExpPattern         = `\A\s*(?:export\s+)?([\w\.]+)(?:\s*=\s*|:\s+?)('(?:\'|[^'])*'|"(?:\"|[^"])*"|[^#\n]+)?\s*(?:\s*\#.*)?\z`
	variableRegExpPattern     = `(\\)?(\$)(\{?([A-Z0-9_]+)?\}?)`
	removeQuotesRegExpPattern = `\A(['"])(.*)(['"])\z`
)

var (
	lineRegExp         = regexp.MustCompile(lineRegExpPattern)
	variableRegExp     = regexp.MustCompile(variableRegExpPattern)
	removeQuotesRegExp = regexp.MustCompile(removeQuotesRegExpPattern)
)

// EnvVar var.
type EnvVar struct {
	Name  string
	Value interface{}
}

// Envs environment variables.
type Envs []EnvVar

func str2Type(s string) interface{} {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	if v, err := strconv.ParseUint(s, 10, 64); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v
	}
	if v, err := strconv.ParseBool(s); err == nil {
		return v
	}
	if v, err := time.ParseDuration(s); err == nil {
		return v
	}
	if strings.Contains(s, "://") {
		if v, err := url.Parse(s); err == nil {
			return v
		}
	}
	return s
}

// Parser parser for env file.
type Parser struct {
}

// Parse env file.
func (l *Parser) Parse(filename string) (Envs, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return l.parse(f)
}

func (l *Parser) parse(f *os.File) (Envs, error) {
	var envVars Envs

	scanner := bufio.NewScanner(f)

	checkMap := make(map[string]struct{}, 64)

	for scanner.Scan() {
		line := scanner.Text()
		name, val, err := l.parseLine(line)
		if err != nil {
			return nil, err
		}
		if name == "" {
			continue
		}
		if _, ok := checkMap[name]; ok {
			return nil, fmt.Errorf("Line `%s` has an unset variable", line)
		}

		envVars = append(envVars, EnvVar{Name: name, Value: str2Type(val)})

		checkMap[name] = struct{}{}
	}
	return envVars, nil
}

func (l *Parser) parseLine(line string) (key string, val string, err error) {
	parts := lineRegExp.FindStringSubmatch(line)

	if len(parts) == 0 {
		st := strings.TrimSpace(line)
		if st == "" || strings.HasPrefix(st, "#") {
			return "", "", nil
		}
		return "", "", fmt.Errorf("Line `%s` doesn't match format", line)
	}
	key = parts[1]
	val = parts[2]

	val = removeQuotesRegExp.ReplaceAllString(val, "$2")

	return key, val, nil
}

// New creates a Parser instance.
func New() *Parser {
	return &Parser{}
}
