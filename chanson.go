// Create Json Streams in Go.
// Chanson makes easy to fetch data from a channel and encode it.
// It is not an encoder it self, by default it relies on json.Encoder but its flexible enough to let you use whatever you want.
package chanson

import (
	"encoding/json"
	"io"
)

type Chanson struct {
	w io.Writer
}

// New returns a new json stream.
// The stream will use w for write the output
func New(w io.Writer) Chanson {
	cs := Chanson{w: w}
	return cs
}

// Object will execute the callback inside an object context
// this is: "{" f() "}"
func (cs *Chanson) Object(f func(obj Object)) {
	cs.w.Write([]byte("{"))
	if f != nil {
		f(Object{cs: cs, empty: true})
	}
	cs.w.Write([]byte("}"))
}

// Object will execute the callback inside an array context
// this is: "[" f() "]"
func (cs *Chanson) Array(f func(a Array)) {
	cs.w.Write([]byte("["))
	if f != nil {
		f(Array{cs: cs, empty: true})
	}
	cs.w.Write([]byte("]"))
}

type Object struct {
	cs    *Chanson
	empty bool
}

// Sets an element into the object
func (obj *Object) Set(key string, val interface{}) {
	if !obj.empty {
		obj.cs.w.Write([]byte(","))
	} else {
		obj.empty = false
	}

	quotedKey, _ := json.Marshal(key)

	obj.cs.w.Write(quotedKey)
	obj.cs.w.Write([]byte(":"))
	handleValue(obj.cs, val)
}

type Array struct {
	cs    *Chanson
	empty bool
}

// Pushes an item into the array
func (a *Array) Push(val interface{}) {
	if !a.empty {
		a.cs.w.Write([]byte(","))
	} else {
		a.empty = false
	}

	handleValue(a.cs, val)
}

func handleValue(cs *Chanson, val interface{}) {
	switch t := val.(type) {
	case func(io.Writer):
		t(cs.w)
	case func(*json.Encoder):
		t(json.NewEncoder(cs.w))
	case func(*Chanson):
		t(cs)
	default:
		json.NewEncoder(cs.w).Encode(val)
	}
}