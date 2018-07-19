package goenv

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var Env = NewEnvSet()

var envShowFlag = flag.Bool("env", false, "")

type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() interface{} { return int64(*i) }

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() interface{} { return uint(*i) }

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() interface{} { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() interface{} { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

type urlValue url.URL

func newURLValue(val url.URL, p *url.URL) *urlValue {
	*p = val
	return (*urlValue)(p)
}

func (u *urlValue) Set(val string) error {
	if !strings.HasPrefix(val, "http") {
		val = "http://" + val
	}
	url, err := url.Parse(val)
	if err != nil {
		return err
	}
	*u = urlValue(*url)
	return nil
}

func (u *urlValue) Get() interface{} { return url.URL(*u) }

func (u *urlValue) String() string { return (*url.URL)(u).String() }

type Value interface {
	String() string
	Set(string) error
}

type EnvVar struct {
	Name     string
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}

type EnvVarSet struct {
	formal map[string]*EnvVar
	output io.Writer
}

func StringVar(p *string, name string, value string) {
	Env.Var(newStringValue(value, p), name)
}

func URLVar(p *url.URL, name string, value url.URL) {
	Env.Var(newURLValue(value, p), name)
}

func BoolVar(p *bool, name string, value bool) {
	Env.Var(newBoolValue(value, p), name)
}

func IntVar(p *int, name string, value int) {
	Env.Var(newIntValue(value, p), name)
}

func Int64Var(p *int64, name string, value int64) {
	Env.Var(newInt64Value(value, p), name)
}

func UintVar(p *uint, name string, value uint) {
	Env.Var(newUintValue(value, p), name)
}

func Uint64Var(p *uint64, name string, value uint64) {
	Env.Var(newUint64Value(value, p), name)
}

func Float64Var(p *float64, name string, value float64) {
	Env.Var(newFloat64Value(value, p), name)
}

func DurationVar(p *time.Duration, name string, value time.Duration) {
	Env.Var(newDurationValue(value, p), name)
}

func (e *EnvVarSet) out() io.Writer {
	if e.output == nil {
		return os.Stderr
	}
	return e.output
}

func (e *EnvVarSet) VisitAll(fn func(*EnvVar)) {
	for _, flag := range e.formal {
		fn(flag)
	}
}

func (e *EnvVarSet) Var(value Value, name string) {
	envVar := &EnvVar{name, value, value.String()}
	_, alreadyThere := e.formal[name]
	if alreadyThere {
		msg := fmt.Sprintf("flag redefined: %s", name)
		fmt.Fprintln(e.out(), msg)
		panic(msg)
	}
	if e.formal == nil {
		e.formal = make(map[string]*EnvVar)
	}
	e.formal[name] = envVar
}

func (e *EnvVarSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(e.out(), err)
	return err
}

func (e *EnvVarSet) printEnv() {
	for _, v := range e.formal {
		val := v.DefValue
		if v.Value.String() != "" {
			val = v.Value.String()
		}
		fmt.Printf("%s=\"%s\" \n", v.Name, val)
	}
	os.Exit(0)
}

func (e *EnvVarSet) Parse() error {
	flag.Parse()

	for _, v := range e.formal {
		value, ok := os.LookupEnv(v.Name)
		if !ok {
			value = v.DefValue
		}
		if err := v.Value.Set(value); err != nil {
			return e.failf("invalid value %q for env %s: %v", value, v.Name, err)
		}
	}

	if *envShowFlag {
		e.printEnv()
	}

	return nil
}

func NewEnvSet() *EnvVarSet {
	return &EnvVarSet{}
}

func Parse() {
	err := Env.Parse()
	if err != nil {
		os.Exit(2)
	}
}
