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
		AvgFloat64(sample)
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

func BenchmarkThanSumPastUint64(b *testing.B) {
	var sample = []uint64{10, 11, 12, 53}
	for i := 0; i < b.N; i++ {
		ThanSumPastUint64(sample, 3, true)
	}
}

func BenchmarkThanSumPast(b *testing.B) {
	var sample = []float64{10.1, 11.1, 12.1, 53.1}
	for i := 0; i < b.N; i++ {
		thanSumPast(sample, true)
	}
}

func TestSum(t *testing.T) {
	var sample1 = []float64{1.1, 2.2, 3.3}
	if SumFloat64(sample1) != 6.6 {
		t.Error("Should be 6.6")
	}
	var sample2 = []uint64{1, 2, 3}
	if SumUint64(sample2) != 6 {
		t.Error("Should be 6")
	}
}

func TestThanPast(t *testing.T) {
	var sample1 = []float64{10.1, 11.1, 12.1, 13.1}
	if !ThanPastFloat64(sample1, 3, true) {
		t.Error("Should be `true`")
	}
	var sample2 = []float64{10.1, 11.1, 12.1, 9.1}
	if !ThanPastFloat64(sample2, 3, false) {
		t.Error("Should be `true`")
	}
	if ThanPastFloat64(sample2, 3, true) {
		t.Error("Should be `false`")
	}
	var sample3 = []uint64{10, 11, 12, 13}
	if !ThanPastUint64(sample3, 3, true) {
		t.Error("Should be `true`")
	}
	var sample4 = []uint64{10, 11, 12, 9}
	if !ThanPastUint64(sample4, 3, false) {
		t.Error("Should be `true`")
	}
}

func TestThanSumPast(t *testing.T) {
	var sample1 = []float64{10.1, 11.1, 12.1, 53.1}
	if !ThanSumPastFloat64(sample1, 3, true) {
		t.Error("Should be `true`")
	}
	var sample2 = []float64{10.1, 11.1, 12.1, 10.1}
	if !ThanSumPastFloat64(sample2, 3, false) {
		t.Error("Should be `true`")
	}
	var sample3 = []uint64{10, 11, 12, 53}
	if !ThanSumPastUint64(sample3, 3, true) {
		t.Error("Should be `true`")
	}
	var sample4 = []uint64{10, 11, 12, 10}
	if !ThanSumPastUint64(sample4, 3, false) {
		t.Error("Should be `true`")
	}
}
