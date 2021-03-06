package statea

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

/* sample kolmogorov values from scipy

for i in range(100):
    v = i / 50.
    print "{%f, %f}," % (v, scipy.special.kolmogorov(v)),

*/
var kolmogorov_values = []struct {
	in  float64
	out float64
}{
	{0.000000, 1.000000}, {0.020000, 1.000000}, {0.040000, 1.000000},
	{0.060000, 1.000000}, {0.080000, 1.000000}, {0.100000, 1.000000},
	{0.120000, 1.000000}, {0.140000, 1.000000}, {0.160000, 1.000000},
	{0.180000, 1.000000}, {0.200000, 1.000000}, {0.220000, 1.000000},
	{0.240000, 1.000000}, {0.260000, 1.000000}, {0.280000, 0.999999},
	{0.300000, 0.999991}, {0.320000, 0.999954}, {0.340000, 0.999829},
	{0.360000, 0.999489}, {0.380000, 0.998715}, {0.400000, 0.997192},
	{0.420000, 0.994524}, {0.440000, 0.990270}, {0.460000, 0.983995},
	{0.480000, 0.975318}, {0.500000, 0.963945}, {0.520000, 0.949694},
	{0.540000, 0.932503}, {0.560000, 0.912423}, {0.580000, 0.889606},
	{0.600000, 0.864283}, {0.620000, 0.836745}, {0.640000, 0.807323},
	{0.660000, 0.776363}, {0.680000, 0.744220}, {0.700000, 0.711235},
	{0.720000, 0.677735}, {0.740000, 0.644019}, {0.760000, 0.610360},
	{0.780000, 0.576998}, {0.800000, 0.544142}, {0.820000, 0.511972},
	{0.840000, 0.480635}, {0.860000, 0.450255}, {0.880000, 0.420929},
	{0.900000, 0.392731}, {0.920000, 0.365715}, {0.940000, 0.339919},
	{0.960000, 0.315364}, {0.980000, 0.292059}, {1.000000, 0.270000},
	{1.020000, 0.249175}, {1.040000, 0.229564}, {1.060000, 0.211140},
	{1.080000, 0.193870}, {1.100000, 0.177718}, {1.120000, 0.162644},
	{1.140000, 0.148605}, {1.160000, 0.135557}, {1.180000, 0.123454},
	{1.200000, 0.112250}, {1.220000, 0.101898}, {1.240000, 0.092352},
	{1.260000, 0.083565}, {1.280000, 0.075494}, {1.300000, 0.068092},
	{1.320000, 0.061318}, {1.340000, 0.055129}, {1.360000, 0.049486},
	{1.380000, 0.044349}, {1.400000, 0.039682}, {1.420000, 0.035449},
	{1.440000, 0.031617}, {1.460000, 0.028154}, {1.480000, 0.025031},
	{1.500000, 0.022218}, {1.520000, 0.019690}, {1.540000, 0.017421},
	{1.560000, 0.015390}, {1.580000, 0.013573}, {1.600000, 0.011952},
	{1.620000, 0.010508}, {1.640000, 0.009223}, {1.660000, 0.008083},
	{1.680000, 0.007072}, {1.700000, 0.006177}, {1.720000, 0.005388},
	{1.740000, 0.004691}, {1.760000, 0.004078}, {1.780000, 0.003540},
	{1.800000, 0.003068}, {1.820000, 0.002654}, {1.840000, 0.002293},
	{1.860000, 0.001977}, {1.880000, 0.001703}, {1.900000, 0.001464},
	{1.920000, 0.001256}, {1.940000, 0.001076}, {1.960000, 0.000921},
	{1.980000, 0.000787},
}

var epsilon = 0.001

func TestKolmogorov(t *testing.T) {
	for _, row := range kolmogorov_values {
		result := Kolmogorov(row.in)
		delta := math.Abs(result - row.out)
		if delta > epsilon {
			t.Errorf("Kolmogorov(%f) == %f, should be %f (error %f)",
				row.in, result, row.out, delta)
		}
	}
}

func TestKolmogorovTest(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	size := 3000
	sample := make([]float64, size)

	// target cdf is uniform
	cdf := func(x float64) float64 {
		return x
	}
	// use uniform for input
	for i := 0; i < size; i++ {
		sample[i] = rand.Float64()
	}
	_, pvalue := KolmogorovTest(sample, cdf)

	if pvalue < 0.05 {
		t.Errorf("KolmogorovTest(uniform == uniform) == %f", pvalue)
	}

	// use normal for input
	desiredStdDev := 0.15
	desiredMean := 0.5
	for i := 0; i < size; i++ {
		sample[i] = rand.NormFloat64()*desiredStdDev + desiredMean
	}
	_, pvalue = KolmogorovTest(sample, cdf)

	if pvalue > 0.05 {
		t.Errorf("KolmogorovTest(normal == uniform) == %f", pvalue)
	}

	// use offset uniform for input
	for i := 0; i < size; i++ {
		sample[i] = rand.Float64() * 1.05
	}
	_, pvalue = KolmogorovTest(sample, cdf)

	if pvalue > 0.05 {
		t.Errorf("KolmogorovTest(offset uniform == uniform) == %f", pvalue)
	}
}
