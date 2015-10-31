package chanson

import (
	"encoding/json"
	"io"
)

type Chanson struct {
	w io.Writer
}

func New(w io.Writer) Chanson {
	cs := Chanson{w: w}
	return cs
}

func (cs *Chanson) Object(f func(obj Object)) {
	cs.w.Write([]byte("{"))
	if f != nil {
		f(Object{cs: cs, empty: true})
	}
	cs.w.Write([]byte("}"))
}

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

func (obj *Object) Set(key string, val interface{}) {
	if !obj.empty {
		obj.cs.w.Write([]byte(","))
	} else {
		obj.empty = false
	}

	obj.cs.w.Write([]byte(`"` + key + `":`))
	handleValue(obj.cs, val)
}

type Array struct {
	cs    *Chanson
	empty bool
}

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