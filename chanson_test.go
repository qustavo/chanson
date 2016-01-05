package chanson

import (
	"bytes"
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
	New(buf).Object(func(obj Object) {
		obj.Set("foo", "bar")
		obj.Set("fun", func(w io.Writer) {
			if _, err := w.Write([]byte(`"val"`)); err != nil {
				panic(err)
			}
		})
	})

	assert.Equal(t, trim(`
	{
		"foo": "bar",
		"fun": "val"
	}`), trim(buf.String()))
}

func TestObjectKeyEncoding(t *testing.T) {

	for _, test := range []struct {
		key      string
		expected string
	}{
		{"key\n", `"key\n"`},
		{"key\t", `"key\t"`},
		{"key\b", `"key\b"`},
		{"key\f", `"key\f"`},
		{"key\\", `"key\\"`},
	} {
		buf := bytes.NewBuffer(nil)
		New(buf).Object(func(obj Object) {
			obj.Set(test.key, 0)
		})

		assert.Equal(t, "{"+test.expected+":0}", trim(buf.String()))
	}
}

func TestArrays(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	New(buf).Array(func(a Array) {
		a.Push(func(w io.Writer) {
			if _, err := w.Write([]byte("1")); err != nil {
				panic(err)
			}
		})
		a.Push(2)
		a.Push(3)
	})

	assert.Equal(t, `[1,2,3]`, trim(buf.String()))
}

func TestObjectSetWithDifferentValueTypes(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	New(buf).Object(func(obj Object) {
		obj.Set("id", 10)
		obj.Set("array", func(arr Array) {
			arr.Push("foo")
			arr.Push("bar")
		})
		obj.Set("obj", func(_obj Object) {
			_obj.Set("foo", "bar")
		})
	})

	assert.Equal(t, trim(`
	{
	  "id": 10,
	  "array": ["foo", "bar"],
	  "obj": {"foo":"bar"}
  	}`), trim(buf.String()))
}

func TestWritesNullWhenValueEncodingFails(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	New(buf).Object(func(obj Object) {
		// func(){} will return an error when tried to be json.Encoder#Encode()
		val := func() {}
		obj.Set("key", val)
	})

	assert.Equal(t, `{"key":null}`, trim(buf.String()))
}

func TestWithChannels(t *testing.T) {
	intCh := make(chan int, 5)
	boolCh := make(chan bool, 2)

	go func() {
		boolCh <- true
		for i := 0; i < cap(intCh); i++ {
			intCh <- i
		}
		boolCh <- false

		close(intCh)
		close(boolCh)

	}()

	buf := bytes.NewBuffer(nil)
	New(buf).Object(func(obj Object) {
		obj.Set("int", func(a Array) {
			for n := range intCh {
				a.Push(n)
			}
		})

		obj.Set("bool", func(a Array) {
			for n := range boolCh {
				a.Push(n)
			}
		})
	})

	assert.Equal(t, trim(`
	{
	  "int": [0,1,2,3,4],
	  "bool": [true,false]
  	}`), trim(buf.String()))

}
