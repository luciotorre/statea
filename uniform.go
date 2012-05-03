/*
Package statea provides data structures that calculate summary statistics
over streams of data sets.
*/
package statea

import (
    "math/rand"
)


/*
UniformSample keeps a number of items (a sample) from all its updates.

This is a reservoir sampling implementation as described in:
    Random Sampling with a Reservoir
    JEFFREY SCOTT VITTER
    http://www.cs.umd.edu/~samir/498/vitter.pdf

The idea is an extension of the algorithm to keep one number from a
stream of unknown length, which is easy to derive.

Always keep the first number, x1. This solves it for N=1.
When presented with the second number, x2, you want to keep it 50%
of the times. So if random() < 1/2, keep x2, else, keep x1.
Now we have a stream of length 2 where every number will be the one
selected 1/2 of the time.

Now, generalizing for the induction step, assume you have seen N items and
you have uniformly choosen one from the previous N. For N+1, you want to keep
x(N+1) number only 1/(N+1) of the times.

So just do that, if random() < 1/(N+1) keep it.

The other numbers had all 1/N probability of being choosen. You will keep the
previous selection only in 1-1/(N+1) of the cases. Multiplying the old
probability by the chance of keeping it we get that the new probability is
now 1/(N+1) as desired.

Now extend that to to have a reservoir of more than one number.

*/

type UniformSample struct {
    Count int // the number of items seen
    size int // the sample size
    Sample []float64 // the sampled numbers
}

func NewUniformSample(size int) *UniformSample {
    if size < 1 {
        return nil
    }
    s := new(UniformSample)
    s.size = size
    s.Count = 0
    s.Sample = make([]float64, size)

    return s
}

func (self *UniformSample) Update(value float64) {

    if self.Count >= self.size {
        r := rand.Int() % self.Count
        if r < self.size {
            self.Sample[r] = value
        }
    } else {
        self.Sample[self.Count] = value
    }

    self.Count += 1
}

