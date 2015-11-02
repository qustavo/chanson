# Chanson [![Build Status](https://travis-ci.org/gchaincl/chanson.svg)](https://travis-ci.org/gchaincl/dotsql)
Create Json Streams in Go

# Why
Let's say you want to json-encode a large amount of data from your database and send it through the network.
So you pack the data into a structure and then you call `json.NewEncoder(w).Encode(theData)`, but this has some implications:
_a)_ If the data is to big, you will end up using a massive amount of memory and
_b)_ your client will need to wait all that time before receiving a single bit.

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

