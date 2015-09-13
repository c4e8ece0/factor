// Package factor contains interface for factor representation
// and type conversion
package factor

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Factorer required for stardart data representation for all data
type Factorer interface {
	Factorize(string) (map[string]Viewer, error)
}

// Custom type converter (value, to_type)
type Converter func(interface{}) (interface{}, error)

// Viewer prepare data for different types
type Viewer interface {
	Set(interface{})                     // set start value
	Get() interface{}                    // get start value
	Error()                              // get last error
	HasError()                           // check error status
	String() string                      // make string
	Int() int64                          // make int64
	Float() float64                      // make float64
	Bool() bool                          // make bool
	Bin() uint8                          // binary, i.e. bool ? 1 : 0
	View(Converter) (interface{}, error) // custom converter
}

// Generic make most convertions by itself ;)
type Generic struct {
	val interface{}
	err error
	Viewer
}

func (g *Generic) Set(v interface{})  { g.val = v }
func (g *Generic) Get() interface{}   { return g.val }
func (g *Generic) Error() interface{} { return g.err }
func (g *Generic) HasError() bool     { return g.err != nil }
func (g *Generic) String() string {
	v, e := g.View(ConvertFunc(ToString))
	if e != nil {
		g.err = e
		return ""
	}
	return v.(string)
}
func (g *Generic) Int() int64 {
	v, e := g.View(ConvertFunc(ToInt))
	if e != nil {
		g.err = e
		return 0
	}
	return v.(int64)
}
func (g *Generic) Float() float64 {
	v, e := g.View(ConvertFunc(ToFloat))
	if e != nil {
		g.err = e
		return 0.0
	}
	return v.(float64)
}
func (g *Generic) Bool() bool {
	v, e := g.View(ConvertFunc(ToBool))
	if e != nil {
		g.err = e
		return false
	}
	return v.(bool)
}
func (g *Generic) Bin() uint8 {
	if g.Bool() {
		return uint8(1)
	}
	return uint8(0)
}
func (g *Generic) View(f Converter) (interface{}, error) {
	return f(g.val.(interface{}))
}

// Internal type ids
const (
	ToString = iota
	ToInt
	ToFloat
	ToBool
)

// Internal Custom Converter example
// was AutoConverter

func ConvertFunc(to int) Converter {
	return func(value interface{}) (interface{}, error) {
		v := reflect.ValueOf(value)
		switch v.Kind() {

		case reflect.String:
			return FromStringTo(value.(string), to)

		case reflect.Int64:
			return FromInt64To(value.(int64), to)

		case reflect.Float64:
			return FromFloat64To(value.(float64), to)

		case reflect.Bool:
			return FromBoolTo(value.(bool), to)

		}
		return 0, errors.New("Unknown source type")
	}
}

func FromStringTo(t string, to int) (interface{}, error) {
	switch to {
	case ToString:
		return t, nil
	case ToInt:
		return strconv.ParseInt(t, 10, 64)
	case ToFloat:
		return strconv.ParseFloat(t, 64)
	case ToBool:
		if t != "" {
			return true, nil
		}
		return false, nil
	}
	return 0, errors.New("Unknown target type")
}

func FromInt64To(t int64, to int) (interface{}, error) {
	switch to {
	case ToString:
		return fmt.Sprintf("%d", t), nil
	case ToInt:
		return t, nil
	case ToFloat:
		return float64(t), nil
	case ToBool:
		if t != 0 {
			return true, nil
		}
		return false, nil
	}
	return 0, errors.New("Unknown target type")

}

func FromFloat64To(t float64, to int) (interface{}, error) {
	switch to {
	case ToString:
		return fmt.Sprintf("%f", t), nil
	case ToInt:
		return int64(t), nil
	case ToFloat:
		return t, nil
	case ToBool:
		if math.Abs(t) > 0 {
			return true, nil
		}
		return false, nil
	}
	return 0, errors.New("Unknown target type")

}

func FromBoolTo(t bool, to int) (interface{}, error) {
	switch to {
	case ToString:
		if t {
			return "true_string", nil
		}
		return "", nil

	case ToInt:
		if t {
			return int64(1), nil
		}
		return int64(0), nil

	case ToFloat:
		if t {
			return float64(1), nil
		}
		return float64(0), nil

	case ToBool:
		return t, nil
	}
	return 0, errors.New("Unknown target type")

}
