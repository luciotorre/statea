package statea

import (
	"container/heap"
	"math/rand"
	"math"
	"time"
)

type Item struct {
	value    float64 // The value of the item; arbitrary.
	priority float64 // The priority of the item in the queue.

	// The index is needed by changePriority and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	// To simplify indexing expressions in these methods, we save a copy of the
	// slice object. We could instead write (*pq)[i].
	a := *pq
	n := len(a)
	a = a[0 : n+1]
	item := x.(*Item)
	item.index = n
	a[n] = item
	*pq = a
}

func (pq *PriorityQueue) Pop() interface{} {
	a := *pq
	n := len(a)
	item := a[n-1]
	item.index = -1 // for safety
	*pq = a[0 : n-1]
	return item
}

// update is not used by the example but shows how to take the top item from
// the queue, update its priority and value, and put it back.
func (pq *PriorityQueue) update(value float64, priority float64) {
	item := heap.Pop(pq).(*Item)
	item.value = value
	item.priority = priority
	heap.Push(pq, item)
}

// changePriority is not used by the example but shows how to change the
// priority of an arbitrary item.
func (pq *PriorityQueue) changePriority(item *Item, priority float64) {
	heap.Remove(pq, item.index)
	item.priority = priority
	heap.Push(pq, item)
}

func (pq PriorityQueue) Rescale(scale func(float64) float64) {
	size := len(pq)

	for i := 0; i < size; i++ {
		pq[i].priority = scale(pq[i].priority)
	}
	
}

type WeigthedSample struct {
	Count int           // the number of items seen
	size  int           // the sample size
	pq    PriorityQueue // the sampled numbers
}

func NewWeigthedSample(size int) *WeigthedSample {
	if size < 1 {
		return nil
	}
	s := new(WeigthedSample)
	s.size = size
	s.Count = 0
	s.pq = make(PriorityQueue, 0, size+1)

	return s
}

func (self *WeigthedSample) Sample() []float64 {
	size := self.size
	if self.size > self.Count {
		size = self.Count
	}

	result := make([]float64, size)
	for i := 0; i < size; i++ {
		result[i] = self.pq[i].value
	}
	return result
}

func (self *WeigthedSample) Rescale(scale func(float64) float64) {
	size := self.size
	if self.size > self.Count {
		size = self.Count
	}

	result := make(PriorityQueue, 0, self.size+1)
	for i := 0; i < size; i++ {
		item := heap.Pop(&self.pq).(*Item)
		item.priority = scale(item.priority)
		heap.Push(&result, item)
	}
	self.pq = result
}

func (self *WeigthedSample) Update(value float64, weight float64) {
	priority := weight / rand.Float64()
	if self.Count >= self.size {
		if self.pq[0].priority < priority {
			item := &Item{
				value:    value,
				priority: priority,
			}
			heap.Push(&self.pq, item)
			heap.Pop(&self.pq)
		}
	} else {
		item := &Item{
			value:    value,
			priority: priority,
		}
		heap.Push(&self.pq, item)
	}

	self.Count += 1
}

/* RESCALE EVERY 10 minutes */
var  RESCALE_THRESHOLD float64 = 60 * 10

type ExponentiallyDecayingSample struct {
	last_t float64
	s *WeigthedSample
	alpha float64
}

func Now() float64 {
	return float64(time.Now().UnixNano()) / 1000000000
}

func NewExponentiallyDecayingSample(size int, alpha float64) *ExponentiallyDecayingSample {
	if size < 1 {
		return nil
	}
	s := new(ExponentiallyDecayingSample)
	s.s = NewWeigthedSample(size)
	s.last_t = Now()
	s.alpha = alpha
	return s
}

func (self *ExponentiallyDecayingSample) Sample() []float64 {
	return self.s.Sample()
}

func (self *ExponentiallyDecayingSample) Update(value float64, timestamp float64) {
	if timestamp - self.last_t > RESCALE_THRESHOLD {
		self.s.Rescale(func (value float64) float64 {
			return value * math.Exp(-self.alpha * (timestamp - self.last_t))
		})
		self.last_t = timestamp
	}
	w := math.Exp(self.alpha * (timestamp - self.last_t))
	self.s.Update(value, w)
}


