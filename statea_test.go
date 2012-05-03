package statea

import (
    "testing"
    )

func Test_Update(t *testing.T) {
    s := NewUniformSample(10)
    if (s.Count != 0) {
        t.Fatalf("Uniform Sample count not zero. (got %d)", s.Count)
    }

    for i := 1; i < 21; i++ {
        s.Update(1)
        if (s.Count != i) {
            t.Fatalf("Uniform Sample count not %d. (got %d)", i, s.Count)
        }
    }

    for i := 0; i < 10; i++ {
        if s.Sample[i] != 1 {
           t.Fatalf("Uniform Sample at %d not 1. (got %d)", i, s.Count)
        }
    }
}

