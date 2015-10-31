package chanson

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func trim(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\n', '\t':
			return -1
		default:
			return r
		}
	}, s)
}

func TestObjectKeyVal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	cs := New(buf)

	cs.Object(func(obj Object) {
		obj.Set("foo", "bar")
		obj.Set("fun", func(w io.Writer) {
			w.Write([]byte(`"val"`))
		})
		obj.Set("a\nnewline", "baz")
	})

	assert.Equal(t, trim(`
	{
		"foo": "bar",
		"fun": "val",
		"a\nnewline": "baz"
	}`), trim(buf.String()))
}

func TestArrays(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	cs := New(buf)

	cs.Array(func(a Array) {
		a.Push(func(w io.Writer) {
			w.Write([]byte("1"))
		})
		a.Push(2)
		a.Push(3)
		a.Push(func(enc *json.Encoder) {
			enc.Encode(4)
		})
	})

	assert.Equal(t, `[1,2,3,4]`, trim(buf.String()))
}

func TestObjectWithArray(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	cs := New(buf)

	cs.Object(func(obj Object) {
		obj.Set("id", 10)
		obj.Set("list", func(cs *Chanson) {
			cs.Array(func(a Array) {
				a.Push(1)
				a.Push(2)
			})
		})
	})

	assert.Equal(t, trim(`
	{
	  "id": 10,
	  "list": [1,2]
  	}`), trim(buf.String()))
}

func TestWithChannels(t *testing.T) {
	intCh := make(chan int, 5)
	boolCh := make(chan bool, 2)

	buf := bytes.NewBuffer(nil)
	cs := New(buf)

	go func() {
		boolCh <- true
		for i := 0; i < cap(intCh); i++ {
			intCh <- i
		}
		boolCh <- false

		close(intCh)
		close(boolCh)

	}()

	cs.Object(func(obj Object) {
		obj.Set("int", func(cs *Chanson) {
			cs.Array(func(a Array) {
				for n := range intCh {
					a.Push(n)
				}
			})
		})
		obj.Set("bool", func(cs *Chanson) {
			cs.Array(func(a Array) {
				for n := range boolCh {
					a.Push(n)
				}
			})
		})
	})

	assert.Equal(t, trim(`
	{
	  "int": [0,1,2,3,4],
	  "bool": [true,false]
  	}`), trim(buf.String()))

}