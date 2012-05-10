package statea

import (
	"math"
	"sort"
)

/* A go re implementation of the cephes kolmogorov code */

/* Kolmogorov's limiting distribution of two-sided test.

Returns the probability that sqrt(n) * max deviation > y,
or that max deviation > y/sqrt(n). The approximation is 
useful for the tail of the distribution when n is large.  
(> ~2666)
*/
func Kolmogorov(y float64) float64 {
	if y < 1.1e-16 {
		return 1.0
	}
	x := -2.0 * y * y
	sign := 1.0
	p := 0.0
	r := 1.0
	t := 0.0

	for {
		t = math.Exp(x * r * r)
		p += sign * t
		if t == 0.0 {
			break
		}
		r += 1.0
		sign = -sign
		if (t / p) <= 1.1e-16 {
			break
		}
	}

	return (p + p)
}

/* Type for Cummulative Distribution Functions */
type CDF func(float64) float64

/* The Kolmogorov-Smirnov test for goodness of fit.

This performs a test of the sampled random variable against a given
distribution cdf(x). Under the null hypothesis the two distributions
are identical, G(x)=F(x). The alternative hypothesis is 'two_sided'

The KS test is only valid for continuous distributions and this implementation
uses the asymptotic formula, so it should be used with big samples (N > 2666).

Returns the test statistic D and the p-value.
*/
func KolmogorovTest(in_sample []float64, cdf CDF) (float64, float64) {
	sample := make([]float64, len(in_sample))
	copy(sample, in_sample)
	sort.Float64s(sample)

	kplus := 0.0
	kminus := 0.0
	n := len(sample)

	for i := 0; i < n; i++ {
		v := sample[i]
		kp := ((float64(i) + 1.0) / float64(n)) - cdf(v)
		km := cdf(v) - (float64(i) / float64(n))

		if kp > kplus {
			kplus = kp
		}
		if km > kminus {
			kminus = km
		}
	}

	D := kplus
	if kminus > kplus {
		D = kminus
	}

	return D, Kolmogorov(D * math.Sqrt(float64(n)))
}
