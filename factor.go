// Package factor contains interface for factor representation
// and type conversion
package factor

import (
	"errors"
	"math"
	"reflect"
	"strconv"
)

// Factorer required for stardart data representation for all data
type Factorer interface {
	Factorize(string) (map[string]Valuer, error)
}

// Custom type converter
type Conteverter func(interface{}, interface{}) (interface{}, error)

// Viewer prepare data for different types
type Viewer interface {
	Set(interface{})  // set start value
	Get() interface{} // get start value
	Error()           // get last error
	HasError()        // check error status
	String() string   // make string
	Int() int64       // make int64
	Float() float64   // make float64
	Bool() bool       // make bool
	Bin() uint8       // binary, i.e. bool ? 1 : 0
	View() Converter  // custom converter
}

// Generic make most convertions by itself ;)
type Generic struct {
	val interface{}
	err error
	Viewer
}

func (g *Generic) Set(v interface{})  { g.val = v }
func (g *Generic) Get() interface{}   { return g.val }
func (g *Generic) Error() interface{} { return g.error }
func (g *Generic) HasError() bool     { return err != nil }
func (g *Generic) String() string {
	v, e := g.View(AutoConverter, string)
	if e != nil {
		g.err = e
		return ""
	}
	return v
}
func (g *Generic) Int() int64 {
	v, e := g.View(AutoConverter, int64)
	if e != nil {
		g.err = e
		return 0
	}
	return v
}
func (g *Generic) Float() float64 {
	v, e := g.View(AutoConverter, float64)
	if e != nil {
		g.err = e
		return ""
	}
	return v
}
func (g *Generic) Bool() bool {
	v, e := g.View(AutoConverter, bool)
	if e != nil {
		g.err = e
		return 0
	}
	return v
}
func (g *Generic) Bin() uint8 {
	if g.Bool() {
		return 1
	}
	return 0
}
func (g *Generic) View(f Conteverter) (interface{}, error) {
	return f(g.val)
}

// Internal Custom Converter example
func Converter(value interface{}, to_type interface{}) (interface{}, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {

	case reflect.String:
		return FromString(string(value), to_type)

	case reflect.Int64:
		return FromInt64(int64(value), to_type)

	case reflect.Float64:
		return FromFloat64(float64(value), to_type)

	case reflect.Bool:
		return FromBool(bool(value), to_type)

	}
}

func FromString(t string, to_type interface{}) (interface{}, error) {
	switch to_type {
	case string:
		return t
	case int64:
		return strconv.ParseInt(t, 10, 64)
	case float64:
		return strconv.ParseFloat(t, 64)
	case bool:
		if t != "" {
			return true
		}
		return false
	}
}

func FromInt64(t int64, to_type interface{}) (interface{}, error) {
	switch to_type {
	case string:
		return fmt.Sprintf("%d", t)
	case int64:
		return t
	case float64:
		return float64(t)
	case bool:
		if t != 0 {
			return true
		}
		return false
	}
}

func FromFloat64(t float64, to_type interface{}) (interface{}, error) {
	switch to_type {
	case string:
		return fmt.Sprintf("%f", t)
	case int64:
		return int64(t)
	case float64:
		return t
	case bool:
		if math.Abs(t) > 0 {
			return true
		}
		return false
	}
}

func FromBool(t bool, to_type interface{}) (interface{}, error) {
	switch to_type {
	case string:
		if t {
			return "true_string"
		}
		return ""

	case int64:
		if t {
			return int64(1)
		}
		return int64(0)

	case float64:
		if t {
			return float64(1)
		}
		return float64(0)

	case bool:
		return t
	}
}
