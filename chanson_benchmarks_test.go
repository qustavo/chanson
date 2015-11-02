package chanson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func buildLargeObject(size int) map[string]string {
	size = size * 900000

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

	json.NewEncoder(ioutil.Discard).Encode(lo)

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
