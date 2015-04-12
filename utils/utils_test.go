package utils

import "testing"

func BenchmarkParseDate(t *testing.B) {
	var strdate = "104/04/01"
	for i := 0; i < t.N; i++ {
		ParseDate(strdate)
	}
}

func BenchmarkSumUint64(t *testing.B) {
	var sample = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < t.N; i++ {
		SumUint64(sample)
	}
}

func BenchmarkAvgUint64(t *testing.B) {
	var sample = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < t.N; i++ {
		AvgUint64(sample)
	}
}

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

func BenchmarkThanPast_false(t *testing.B) {
	var sample = []float64{20.2, 20.3, 100.25, 100.75}
	for i := 0; i < t.N; i++ {
		thanPast(sample, false)
	}
}

func BenchmarkThanPastFloat64(t *testing.B) {
	var sample = []float64{20.2, 20.3, 100.25, 100.75}
	for i := 0; i < t.N; i++ {
		ThanPastFloat64(sample, 3, true)
	}
}

func BenchmarkThanPastUint64(t *testing.B) {
	var sample = []uint64{20, 23, 125, 105}
	for i := 0; i < t.N; i++ {
		ThanPastUint64(sample, 3, true)
	}
}
