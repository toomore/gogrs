package utils

import (
	"testing"
	"time"
)

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

func TestRanInt(t *testing.T) {
	if (RandInt() - RandInt() + RandInt() - RandInt()) == 0 {
		t.Error("Should not be the same.")
	}
}

func TestParseDate(t *testing.T) {
	var sample1 = "104/2/28"
	if ParseDate(sample1) != time.Date(2015, 2, 28, 0, 0, 0, 0, TaipeiTimeZone) {
		t.Error("Should be 2015/2/28")
	}
	var sample2 = "104/2/29"
	if ParseDate(sample2) != time.Date(2015, 3, 1, 0, 0, 0, 0, TaipeiTimeZone) {
		t.Error("Should be 2015/3/1")
	}
	var sample3 = "104/4/31"
	if ParseDate(sample3) != time.Date(2015, 5, 1, 0, 0, 0, 0, TaipeiTimeZone) {
		t.Error("Should be 2015/5/1")
	}
}

func TestAvg(t *testing.T) {
	var sample1 = []float64{3.3, 6.6, 9.9}
	if AvgFloat64(sample1) != 6.59 {
		t.Error("Should be 6.59")
	}
	var sample2 = []uint64{3, 6, 9}
	if AvgUint64(sample2) != 6 {
		t.Error("Should be 6")
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

func BenchmarkCountCountine(b *testing.B) {
	var sample1 = []float64{10.1, 11.1, 12.1, 53.1}
	for i := 0; i < b.N; i++ {
		CountCountineFloat64(sample1)
	}
}

func TestCountCountine(t *testing.T) {
	var sample1 = []float64{10.1, 11.1, 12.1, 53.1}  // 4 true
	var sample2 = []float64{10.1, 11.1, -12.1, 53.1} // 1 true
	var sample3 = []float64{10.1, 11.1, 12.1, -53.1} // 1 false
	if times, max := CountCountineFloat64(sample1); times != 4 && max != true {
		t.Error("Should be `4 true`")
	}
	if times, max := CountCountineFloat64(sample2); times != 1 && max != true {
		t.Error("Should be `1 true`")
	}
	if times, max := CountCountineFloat64(sample3); times != 1 && max != false {
		t.Error("Should be `1 false`")
	}
}

func TestCalDiff(t *testing.T) {
	var sampleA = []float64{10.0, 11.1, 12.2, 13.3}
	var sampleB = []float64{12.2, 11.1, 10.0}
	var result = CalDiffFloat64(sampleA, sampleB)
	if result[2] != float64(13.3)-float64(10.0) {
		t.Error("Wrong cal.")
	}
	result = CalDiffFloat64(sampleB, sampleA)
	if result[2] != float64(10.0)-float64(13.3) {
		t.Error("Wrong cal.")
	}

	var sampleC = []int64{10, 11, 12, 13}
	var sampleD = []int64{12, 11, 10}
	var result2 = CalDiffInt64(sampleC, sampleD)
	if result2[2] != 13-10 {
		t.Error("Wrong cal.")
	}
	result2 = CalDiffInt64(sampleD, sampleC)
	if result2[2] != 10-13 {
		t.Error("Wrong cal.")
	}
}

func TestDelta(t *testing.T) {
	var sample1 = []float64{10.0, 11.0, 9.0}
	var sample2 = []int64{10, 11, 9}
	t.Log(DeltaFloat64(sample1))
	t.Log(DeltaInt64(sample2))
}

func BenchmarkDeltafloat64(b *testing.B) {
	var sample = []float64{10.0, 11.0, 9.0}
	for i := 0; i < b.N; i++ {
		DeltaFloat64(sample)
	}
}

func BenchmarkDeltaInt64(b *testing.B) {
	var sample = []int64{10, 11, 9}
	for i := 0; i < b.N; i++ {
		DeltaInt64(sample)
	}
}

func BenchmarkCalDiff(b *testing.B) {
	var sampleA = []float64{10.0, 11.1, 12.2, 13.3}
	var sampleB = []float64{12.2, 11.1, 10.0}
	for i := 0; i < b.N; i++ {
		CalDiffFloat64(sampleA, sampleB)
	}
}

func TestSD(t *testing.T) {
	var sample = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if SDUint64(sample) != 2.8722813232690143 {
		t.Error("Should be 2.8722813232690143")
	}
}

func BenchmarkSD(b *testing.B) {
	var sample = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		SDUint64(sample)
	}
}
