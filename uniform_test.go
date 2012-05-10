package statea

import (
	"testing"
)

func Test_R_Update(t *testing.T) {
	s := NewUniformSampleR(10)
	if s.Count != 0 {
		t.Fatalf("Uniform Sample count not zero. (got %d)", s.Count)
	}

	for i := 1; i < 21; i++ {
		s.Update(1)
		if s.Count != i {
			t.Fatalf("Uniform Sample count not %d. (got %d)", i, s.Count)
		}
	}

	for i := 0; i < 10; i++ {
		if s.Sample[i] != 1 {
			t.Fatalf("Uniform Sample at %d not 1. (got %d)", i, s.Count)
		}
	}
}

/* Use Kolmogorov-Smirnov test to validate that the sample is random */

func Test_R_Kolmogorov(t *testing.T) {
	start := 0
	end := 100000
	cdf := func(x float64) float64 {
		return (x - float64(start)) / float64(end-start)
	}

	s := NewUniformSampleR(10000)

	for i := start; i < end; i++ {
		s.Update(float64(i))
	}

	D, pvalue := KolmogorovTest(s.Sample, cdf)

	if pvalue < 0.005 {
		t.Errorf("KolmogorovTest(sample) == %f, D == %f", pvalue, D)
	}

}
