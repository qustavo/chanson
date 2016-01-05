package chanson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func buildLargeObject(size int) map[string]string {
	obj := make(map[string]string, size)

	for i := 0; i < size; i++ {
		key := fmt.Sprintf("key-%d", i)
		val := fmt.Sprintf("val-%d", i)
		obj[key] = val
	}
	return obj
}

func BenchmarkNativeEncoder(b *testing.B) {
	lo := buildLargeObject(b.N)
	b.ResetTimer()

	err := json.NewEncoder(ioutil.Discard).Encode(lo)
	if err != nil {
		panic(err)
	}

}

func BenchmarkChanson(b *testing.B) {
	lo := buildLargeObject(b.N)
	b.ResetTimer()

	New(ioutil.Discard).Object(func(obj Object) {
		for k, v := range lo {
			obj.Set(k, v)
		}
	})
}
