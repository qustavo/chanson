package chanson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// We will read from channels inside a goroutine
// While we channels are open we will stream the output as json
func ExampleChanson() {
	ch := make(chan int)
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		close(ch)
	}()

	buf := bytes.NewBuffer(nil)
	cs := New(buf)
	cs.Array(func(a Array) {
		for i := range ch {
			a.Push(i)
		}
	})

	fmt.Printf("%v", buf.String())
}

func ExampleObject() {
	buf := bytes.NewBuffer(nil)
	cs := New(ioutil.Discard)
	cs.Object(func(obj Object) {
		obj.Set("foo", "bar")
		obj.Set("fun", func(enc *json.Encoder) {
			_ = enc.Encode([]int{1, 2, 3})
		})
	})

	fmt.Printf("%v", buf.String())
}
