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

"Algorithm R (which is is a reservoir algorithm due to Alan Waterman) works
as follows: When the (t + 1)st record in the file is being processed, for t >= n, the
n candidates form a random sample of the first t records. The (t + 1)st record
has a n/(t + 1) chance of being in a random sample of size n of the first t + 1
records, and so it is made a candidate with probability n/(t + 1). The candidate
it replaces is chosen randomly from the n candidates. It is easy to see that the
resulting set of n candidates forms a random sample of the first t + 1 records."

*/

type UniformSampleR struct {
    Count int // the number of items seen
    size int // the sample size
    Sample []float64 // the sampled numbers
}

func NewUniformSampleR(size int) *UniformSampleR {
    if size < 1 {
        return nil
    }
    s := new(UniformSampleR)
    s.size = size
    s.Count = 0
    s.Sample = make([]float64, size)

    return s
}

func (self *UniformSampleR) Update(value float64) {

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

