package statea

import (
	"testing"
	"fmt"
)

func Test_WeigthedSample_Update(t *testing.T) {
	s := NewWeigthedSample(10)
	if s.Count != 0 {
		t.Fatalf("Uniform Sample count not zero. (got %d)", s.Count)
	}

	for i := 1; i < 21; i++ {
		s.Update(1, 1)
		if s.Count != i {
			t.Fatalf("Uniform Sample count not %d. (got %d)", i, s.Count)
		}
	}

	for i := 0; i < 10; i++ {
		if s.Sample()[i] != 1 {
			t.Fatalf("Uniform Sample at %d not 1. (got %d)", i, s.Count)
		}
	}
}

/* Use Kolmogorov-Smirnov test to validate that the sample is random */

func Test_WeigthedSample_Kolmogorov(t *testing.T) {
	end := 1000
	cdf := func(x float64) float64 {
		px := x / float64(end)
		return px * px
	}

	s := NewWeigthedSample(100)

	for i := 0; i < end; i++ {
		s.Update(float64(i), float64(i))
	}

	sample := s.Sample()
	D, pvalue := KolmogorovTest(sample, cdf)

	if pvalue < 0.005 {
		t.Errorf("KolmogorovTest(sample) == %f, D == %f", pvalue, D)
	}

}

func Test_WeigthedSample_Rescale_Kolmogorov(t *testing.T) {
	end := 1000
	cdf := func(x float64) float64 {
		px := x / float64(end)
		return px * px
	}

	s := NewWeigthedSample(100)

	for i := 0; i < 500; i++ {
		s.Update(float64(i), float64(i))
	}
	s.Rescale(func(weight float64) float64 {
		return weight * 2
	})
	for i := 0; i < 500; i++ {
		s.Update(500.0 + float64(i), 2 * (500 + float64(i)))
	}

	sample := s.Sample()
	D, pvalue := KolmogorovTest(sample, cdf)

	if pvalue < 0.005 {
		t.Errorf("KolmogorovTest(sample) == %f, D == %f", pvalue, D)
	}

}

func Benchmark_WeigthedSample(t *testing.B) {
	end := 1000000
	s := NewWeigthedSample(1024)

	for i := 0; i < end; i++ {
		s.Update(float64(i), float64(i))
	}
}

func Benchmark_ExponentiallyDecayingSample(t *testing.B) {
	end := 1000000
	
	ts := Now()
	s := NewExponentiallyDecayingSample(1024, 0.3)
	start := Now()
	
	for i := 0; i < end; i++ {
		s.Update(float64(i), start + float64(i) / 10000000)
	}
	fmt.Printf("duration: %f\n", Now() - ts)
}
