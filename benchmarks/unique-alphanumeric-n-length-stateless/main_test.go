package unique_alphanumeric_n_length_stateless

import (
	"fmt"
	"testing"
)

func BenchmarkGenerators(b *testing.B) {
	tests := []struct {
		name string
		fn   func() string
	}{
		{"MD5Timestamp", generateMD5Timestamp},
		{"Random", generateRandom},
		{"UUIDv4", generateUUIDv4},
		{"generateWithRandEntropy", generateWithRandEntropy},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tt.fn()
			}
		})
	}
}

func TestCollisionRate(t *testing.T) {
	tests := []struct {
		name string
		fn   func() string
	}{
		{"MD5Timestamp", generateMD5Timestamp},
		{"Random", generateRandom},
		{"UUIDv4", generateUUIDv4},
		{"generateWithRandEntropy", generateWithRandEntropy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := make(map[string]struct{})
			collision := 0
			for i := 0; i < sampleSize; i++ {
				val := tt.fn()
				if _, exists := set[val]; exists {
					collision++
				} else {
					set[val] = struct{}{}
				}
			}
			fmt.Printf("%s: %d collisions out of %d samples\n", tt.name, collision, sampleSize)
		})
	}
}
