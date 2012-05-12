package statea

import (
	"testing"
	"fmt"
	"math"
	"math/rand"
	"time"
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
		m := s.pq[0].priority
		for n := range s.pq {
			if s.pq[n].priority < m {
				// we depend on this to do peek
				t.Errorf("Minimum not in first position.")
			}
		}
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

func Test_ExponentiallyDecayingSample_Kolmogorov(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	end := 10000000
	scale := 100000.0

	s := NewExponentiallyDecayingSample(1024, 0.3)
	s.last_t = 0
	
	for i := 0; i < end; i++ {
		s.Update(float64(i)/scale, float64(i)/scale)
	}

	cdf := func(x float64) float64 {
		w := math.Exp(0.3 * (x - s.last_t))
		top_w := math.Exp(0.3 * (float64(end-1)/scale - s.last_t))
		return w / top_w
	}
	sample := s.Sample()
	D, pvalue := KolmogorovTest(sample, cdf)

	if pvalue < 0.005 {
		t.Errorf("KolmogorovTest(sample) == %f, D == %f", pvalue, D)
	}
	
}

func Benchmark_ExponentiallyDecayingSample(t *testing.B) {
	end := 100000
	now := Now()
	s := NewExponentiallyDecayingSample(1024, 0.3)
	start := Now()
	
	for i := 0; i < end; i++ {
		s.Update(float64(i), start + float64(i) / 10000000)
	}
	fmt.Printf("duration: %f\n", Now() - now)
}
