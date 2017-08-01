package goenv

import (
	"testing"
	"time"
)

type ts struct {
	TestBool     bool
	TestInt      int
	TestInt64    int64
	TestUint     uint
	TestUint64   uint64
	TestString   string
	TestFloat64  float64
	TestDuration time.Duration
}

func TestDefault(t *testing.T) {
	d, _ := time.ParseDuration("2h45m")

	ts := &ts{}
	BoolVar(&ts.TestBool, "TEST_BOOL", false)
	IntVar(&ts.TestInt, "TEST_INT", 0)
	Int64Var(&ts.TestInt64, "TEST_INT64", 1)
	UintVar(&ts.TestUint, "TEST_UINT", 2)
	Uint64Var(&ts.TestUint64, "TEST_UINT64", 3)
	StringVar(&ts.TestString, "TEST_STRING", "0")
	Float64Var(&ts.TestFloat64, "TEST_FLOAT64", 4.25)
	DurationVar(&ts.TestDuration, "TEST_DURATION", d)

	Parse()

	if ts.TestBool != false {
		t.Error("TestBool should be false")
	}
	if ts.TestInt != 0 {
		t.Error("TestInt should be 0")
	}
	if ts.TestInt64 != 1 {
		t.Error("TestInt64 should be 1")
	}
	if ts.TestUint64 != 3 {
		t.Error("TestBool should be false")
	}
	if ts.TestString != "0" {
		t.Error("TestBool should be false")
	}
	if ts.TestFloat64 != 4.25 {
		t.Error("TestBool should be false")
	}
	if ts.TestDuration != d {
		t.Error("TestBool should be false")
	}
}

func TestNewBoolValue(t *testing.T) {
	var p bool

	v := newBoolValue(true, &p)
	if v.Get().(bool) != true {
		t.Error("Get newBoolValue should be true")
	}
	v.Set("false")
	if p != false {
		t.Error("Get newBoolValue should be false")
	}
}

func TestNewIntValue(t *testing.T) {
	var p int

	v := newIntValue(100, &p)
	if v.Get().(int) != 100 {
		t.Error("Get newIntValue should be 100")
	}
	v.Set("50")
	if p != 50 {
		t.Error("Get newIntValue should be 50")
	}
}

func TestNewInt64Value(t *testing.T) {
	var p int64

	v := newInt64Value(100, &p)
	if v.Get().(int64) != 100 {
		t.Error("Get newInt64Value should be 100")
	}
	v.Set("50")
	if p != 50 {
		t.Error("Get newInt64Value should be 50")
	}
}

func TestNewUintValue(t *testing.T) {
	var p uint

	v := newUintValue(100, &p)
	if v.Get().(uint) != 100 {
		t.Error("Get newUintValue should be 100")
	}
	v.Set("50")
	if p != 50 {
		t.Error("Get newUintValue should be 50")
	}
}

func TestNewUint64Value(t *testing.T) {
	var p uint64

	v := newUint64Value(100, &p)
	if v.Get().(uint64) != 100 {
		t.Error("Get newUint64Value should be 100")
	}
	v.Set("50")
	if p != 50 {
		t.Error("Get newUint64Value should be 50")
	}
}

func TestNewFloat64Value(t *testing.T) {
	var p float64

	v := newFloat64Value(100.999, &p)
	if v.Get().(float64) != 100.999 {
		t.Error("Get newFloat64Value should be 100.999")
	}
	v.Set("50.999")
	if p != 50.999 {
		t.Error("Get newFloat64Value should be 50.999")
	}
}

func TestNewDurationValue(t *testing.T) {
	var p time.Duration

	v := newDurationValue(time.Hour, &p)
	if v.Get().(time.Duration) != time.Hour {
		t.Error("Get newDurationValue should be 1h")
	}
	v.Set("2h")
	if p != time.Hour*2 {
		t.Error("Get newFloat64Value should be 2h")
	}
}

func TestNewStringValue(t *testing.T) {
	var p string

	v := newStringValue("hello", &p)
	if v.Get().(string) != "hello" {
		t.Error("Get newFloat64Value should be hello")
	}
	v.Set("hello world")
	if p != "hello world" {
		t.Error("Get newFloat64Value should be hello world")
	}
}
