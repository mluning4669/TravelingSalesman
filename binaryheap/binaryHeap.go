package binaryheap

import (
	"Algorithms/Assignment_2/Golang/graphs"
	"errors"
)

//Heap a heap for implmenting a priority queue
type Heap struct {
	Arr      []*graphs.Node
	size     int
	Capacity int
	dict     map[*graphs.Node]int //position in the text
}

//StartHeap Initializes a heap of size N
func StartHeap(n int) *Heap {
	h := Heap{size: 0, Capacity: n + 1} //using 1 based indexing to make the math more managable
	h.Arr = make([]*graphs.Node, n+1)
	h.dict = make(map[*graphs.Node]int)

	return &h
}

//ExtractMin removes minimum element of heap
func (h *Heap) ExtractMin() *graphs.Node {
	min := h.Arr[1]
	h.Arr[1] = h.Arr[h.size]
	h.Arr[h.size] = nil
	h.heapifyDown(1)
	h.size--
	return min
}

//Delete removes element at location i in heap array. Also acts as Delete(elem) when you pass Delete(dict[elem])
func (h *Heap) Delete(i int) {
	h.Arr[i] = h.Arr[h.size]
	h.Arr[h.size] = nil
	h.heapifyDown(i)
	h.size--
}

//Insert insert an element into the heap
func (h *Heap) Insert(elem *graphs.Node) error {
	if h.size+1 >= h.Capacity {
		return errors.New("Heap at capacity")
	}

	h.size++
	h.Arr[h.size] = elem  //insert element at element size
	h.dict[elem] = h.size //store element and size in position dictionary
	h.heapifyUp(h.size)   //put heap in heap order

	return nil
}

//FindMin finds minimum element in heap but doesn't remove it
func (h *Heap) FindMin() *graphs.Node {
	return h.Arr[1]
}

//ChangeKey change the index of current by inserting an deleting as needed
func (h *Heap) ChangeKey(current *graphs.Node, newKey int) error {
	key, prs := h.dict[current]
	if !prs {
		return errors.New("Item not found")
	}

	if newKey > h.Capacity {
		return errors.New("Key excedes heap capacity")
	}

	if newKey > h.size+1 {
		return errors.New("Changing key would violate heap shape (key > size + 1)")
	}

	//if newKey is at the end of the heap array then simply insert current and then delete
	if newKey == h.size+1 {
		h.Insert(current)
		h.Delete(key)
		return nil
	}

	if h.Arr[newKey] != nil {
		//swap old element at newKey with current
		h.swap(key, newKey)
		//heapify up at newKey
		h.heapifyUp(newKey)
		//delete old element
		h.Delete(key)
	}

	return nil
}

func (h *Heap) swap(i int, j int) {
	h.dict[h.Arr[i]] = j //update dictionary
	h.dict[h.Arr[j]] = i

	temp := h.Arr[j]
	h.Arr[j] = h.Arr[i]
	h.Arr[i] = temp
}

func (h *Heap) heapifyUp(i int) {
	if i > 1 {
		j := i / 2 //this works thanks to integer division

		if *h.Arr[j].Weight > *h.Arr[i].Weight {
			h.swap(i, j)

			h.heapifyUp(j)
		}
	}
}

func (h *Heap) heapifyDown(i int) {
	var j = i

	left := 2 * i
	right := 2*i + 1

	if left < h.size && *h.Arr[j].Weight > *h.Arr[left].Weight {
		j = left
	}
	if right < h.size && *h.Arr[j].Weight > *h.Arr[right].Weight {
		j = right
	}

	if j != i {
		h.swap(i, j)

		h.heapifyDown(j)
	}
}
