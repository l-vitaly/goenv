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

// EnvSet environment variables.
type EnvSet map[string]EnvVar

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
	if v, err := url.Parse(s); err == nil && v.Scheme != "" {
		return v
	}
	return s
}

// Parser parser for env file.
type Parser struct {
}

// Parse env file.
func (l *Parser) Parse(filename string) (EnvSet, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return l.parse(f)
}

func (l *Parser) parse(f *os.File) (EnvSet, error) {
	envVars := make(EnvSet, 32)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		name, val, err := l.parseLine(line)
		if err != nil {
			return nil, err
		}
		if name == "" {
			continue
		}
		if _, ok := envVars[name]; ok {
			return nil, fmt.Errorf("Line `%s` has an unset variable", line)
		}
		envVars[name] = EnvVar{Name: name, Value: str2Type(val)}
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
