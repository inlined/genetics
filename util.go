package genetics

type tie struct {
	index   int
	fitness Fitness
}

type maxTieHeap []tie

func (h maxTieHeap) Len() int           { return len(h) }
func (h maxTieHeap) Less(i, j int) bool { return h[i].fitness > h[j].fitness }
func (h maxTieHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push is unsupported in this pacakge
func (h maxTieHeap) Push(x interface{}) {
	panic("maxTieHeap.Push() unsupported")
}

// Pop unsupported in this package
func (h maxTieHeap) Pop() interface{} {
	panic("maxTieHeap.Pop() unsupported")
}
