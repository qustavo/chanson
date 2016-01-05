// Package chanson provides a flexible way to construct JSON documents.
// As chanson populates Arrays and Objects from functions, it's perfectly suitable for streaming jsons as you build it (see examples).
// It is not an encoder it self, by default it relies on json.Encoder but its flexible enough to let you use whatever you want.
package chanson

import (
	"encoding/json"
	"io"
	"strconv"
)

// Chanson struct is the handler representing the current json being encoded
type Chanson struct {
	w   io.Writer
	enc *json.Encoder
}

// Value is the types that functions like Array.Push() or Object.Set() can accepts as values.
// Custom Value types are:
//	- func(Array)
//	- func(Object)
//	- func(io.Writer)
// If Value is none of the above, it will be encoded using json.Encoder
type Value interface{}

// New returns a new json stream.
// The stream will use w for write the output
func New(w io.Writer) Chanson {
	cs := Chanson{
		w:   w,
		enc: json.NewEncoder(w),
	}
	return cs
}

// Object will execute the callback inside an object context
// this is: "{" f() "}"
func (cs Chanson) Object(f func(Object)) {
	_, _ = cs.w.Write([]byte("{"))
	if f != nil {
		f(Object{cs: &cs, empty: true})
	}
	_, _ = cs.w.Write([]byte("}"))
}

// Array will execute the callback inside an array context
// this is: "[" f() "]"
func (cs Chanson) Array(f func(Array)) {
	_, _ = cs.w.Write([]byte("["))
	if f != nil {
		f(newArray(&cs))
	}
	_, _ = cs.w.Write([]byte("]"))
}

// Object struct represent a Json Object ({}).
type Object struct {
	cs    *Chanson
	empty bool
}

// Set add an element into the object
func (obj *Object) Set(key string, val Value) {
	if !obj.empty {
		_, _ = obj.cs.w.Write([]byte(","))
	} else {
		obj.empty = false
	}

	_, _ = obj.cs.w.Write([]byte(strconv.Quote(key)))
	_, _ = obj.cs.w.Write([]byte(":"))
	handleValue(*obj.cs, val)
}

// Array struct represents a Json Array ([]).
type Array struct {
	cs    *Chanson
	empty bool
}

func newArray(cs *Chanson) Array {
	return Array{cs: cs, empty: true}
}

// Push pushes an item into the array
func (a *Array) Push(val Value) {
	if !a.empty {
		_, _ = a.cs.w.Write([]byte(","))
	} else {
		a.empty = false
	}

	handleValue(*a.cs, val)
}

func handleValue(cs Chanson, val Value) {
	switch t := val.(type) {
	case func(Array):
		cs.Array(t)
	case func(Object):
		cs.Object(t)
	case func(io.Writer):
		t(cs.w)
	default:
		err := cs.enc.Encode(val)
		if err != nil {
			//TODO: should panic?!
			_, _ = cs.w.Write([]byte("null"))
		}
	}
}
