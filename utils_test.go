package gogrs

import "testing"

func BenchmarkRanInt(t *testing.B) {
	for i := 0; i < t.N; i++ {
		RandInt()
	}
}
func BenchmarkRanCryptInt(t *testing.B) {
	for i := 0; i < t.N; i++ {
		RandCryptInt()
	}
}
