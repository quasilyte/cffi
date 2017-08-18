package funcall

import (
	"testing"
)

func BenchmarkDirect(b *testing.B) {
	xs := []int{1, 2, 3}
	for i := 0; i < b.N; i++ {
		counter += cgoadd(1, 1)
		counter += cgomystrlen("123")
		counter += len(cgofoo())
		counter += cgosum(xs)
	}
}

func BenchmarkWrapped(b *testing.B) {
	xs := []int{1, 2, 3}
	for i := 0; i < b.N; i++ {
		counter += add(1, 1)
		counter += mystrlen("123")
		counter += len(foo())
		counter += sum(xs)
	}
}
