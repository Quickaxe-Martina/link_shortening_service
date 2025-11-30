package main

import "testing"

func BenchmarkDummy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = "hello" + "world"
	}
}
