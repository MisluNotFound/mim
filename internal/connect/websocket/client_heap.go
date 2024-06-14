package websocket

type clientHeap []*Client

func (h clientHeap) Len() int { return len(h) }

func (h clientHeap) Less(i, j int) bool {
	return h[i].HeartBeat.Before(h[j].HeartBeat)
}

func (h clientHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *clientHeap) Push(x interface{}) {
	*h = append(*h, x.(*Client))
}

func (h *clientHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
