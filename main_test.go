package main

import (
	"testing"
)

type S struct {
	data string
}

// Value receiver
func (s S) Read() string {
	return s.data
}

// Pointer receiver
func (s *S) Read2() string {
	return s.data
}

func BenchmarkValueReceiver(b *testing.B) {
	s := S{data: "test"}
	for i := 0; i < b.N; i++ {
		s.Read()
	}
}

func BenchmarkPointerReceiver(b *testing.B) {
	s := S{data: "test"}
	for i := 0; i < b.N; i++ {
		(&s).Read2()
	}
}
