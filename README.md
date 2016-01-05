# Chanson [![Build Status](https://travis-ci.org/gchaincl/chanson.svg)](https://travis-ci.org/gchaincl/chanson) [![Coverage Status](https://coveralls.io/repos/gchaincl/chanson/badge.svg?branch=coveralls&service=github)](https://coveralls.io/github/gchaincl/chanson?branch=coveralls)
Package chanson provides a flexible way to construct JSON documents.
As chanson populates Arrays and Objects from functions, it's perfectly suitable for streaming jsons as you build it.
It is not an encoder it self, by default it relies on json.Encoder but its flexible enough to let you use whatever you want.

# Example

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/gchaincl/chanson"
)

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		close(ch)
	}()

	buf := bytes.NewBuffer(nil)
	cs := chanson.New(buf)
	cs.Array(func(a chanson.Array) {
		for i := range ch {
			a.Push(i)
		}
	})

	fmt.Printf("%v", buf.String())
}
```

For more examples and documentarion see the [Godoc](http://godoc.org/github.com/gchaincl/chanson).

