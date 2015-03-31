package gogrs

import "testing"

func BenchmarkRanInt(t *testing.B) {
	for i := 0; i < t.N; i++ {
		RandInt()
	}
}

func BenchmarkSumFloat64(t *testing.B) {
	var sample = []float64{20.2, 20.3, 100.25, 100.75}
	for i := 0; i < t.N; i++ {
		SumFloat64(sample)
	}
}

func BenchmarkAvgFloat64(t *testing.B) {
	var sample = []float64{20.2, 20.3, 100.25, 100.75}
	for i := 0; i < t.N; i++ {
		AvgFlast64(sample)
	}
}

func BenchmarkThanPast(t *testing.B) {
	var sample = []float64{20.2, 20.3, 100.25, 100.75}
	for i := 0; i < t.N; i++ {
		thanPast(sample, true)
	}
}

func BenchmarkThanPastFloat64(t *testing.B) {
	var sample = []float64{20.2, 20.3, 100.25, 100.75}
	for i := 0; i < t.N; i++ {
		ThanPastFloat64(sample, 3, true)
	}
}
