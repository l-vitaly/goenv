package goenv

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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

// EnvVars environment variables.
type EnvVars map[string]string

// Loader loading env file.
type Loader struct {
}

// Load load nev file.
func (l *Loader) Load(filename string) (EnvVars, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return l.parse(f)
}

func (l *Loader) parse(f *os.File) (EnvVars, error) {
	envVars := make(EnvVars, 16)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		key, val, err := l.parseLine(line)
		if err != nil {
			return nil, err
		}
		if key == "" {
			continue
		}
		if _, ok := envVars[key]; ok {
			return nil, fmt.Errorf("Line `%s` has an unset variable", line)
		}
		envVars[key] = val
	}
	return envVars, nil
}

func (l *Loader) parseLine(line string) (key string, val string, err error) {
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
