package timing_wheel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

func TestLevel(t *testing.T) {
	lvl, lvl_idx, idx := location(232)
	t.Logf("location : %d, idx in location : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = location(256)
	t.Logf("location : %d, idx in location : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = location(635)
	t.Logf("location : %d, idx in location : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = location(1048790)
	t.Logf("location : %d, idx in location : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = location(4194451)
	t.Logf("location : %d, idx in location : %d, idx: %d", lvl, lvl_idx, idx)

	lvl, lvl_idx, idx = location(256*64*64*64 + 256*64 + 123)
	t.Logf("location[256*64*64*64 + 256*64 + 123] -  lvl:%d, idx in location : %d, idx: %d", lvl, lvl_idx, idx)
	round := 256
	for i := 1; i < lvl; i++ {
		round *= 64
	}
	lvl, lvl_idx, idx = location(256*64*64*64 + 256*64 + 123 - uint64(round))
	t.Logf("location[%d = (256*64*64*64 + 256*64 + 123 - %d)] -  lvl:%d, idx in location : %d, idx: %d", 256*64*64*64+256*64+123-uint64(round), round, lvl, lvl_idx, idx)
}

func TestPointer(t *testing.T) {
	/*empty := &Node{
		index: -1,
	}*/
	node := &Node{
		index: 1,
		do:    nil,
		next:  nil,
		//prev:  nil,
	}
	nextNode := &Node{
		index: 10,
	}
	pointer := (*unsafe.Pointer)(unsafe.Pointer(&node.next))
	t.Logf("next pointer : %v", pointer)
	//t.Logf("equals : %v", pointer == unsafe.Pointer(empty))
	swapped := atomic.CompareAndSwapPointer(pointer, nil, unsafe.Pointer(nextNode))
	t.Logf("swap : %v", swapped)
	if swapped {
		t.Logf("new next node : %v", *node.next)
	}
	//atomic.StorePointer(&pointer, unsafe.Pointer(nextNode))
	//t.Logf("new next node : %v", *node.next)

}

func TestSlot_Append(t *testing.T) {
	slot := &Slot{
		root:    nil,
		current: nil,
	}
	starting := sync.WaitGroup{}
	starting.Add(1)
	all := sync.WaitGroup{}
	for i := 1; i <= 100; i++ {
		node := &Node{
			index:  uint64(i),
			do:     nil,
			depth:  0,
			parent: nil,
			left:   nil,
			right:  nil,
			next:   nil,
		}
		node.do = func() {
			t.Logf("index[%v], depth[%v], next[%v]", node.index, node.depth, *node.next)
		}
		all.Add(1)
		go func() {
			starting.Wait()
			slot.Append(node)
			all.Done()
		}()
	}
	starting.Done()
	all.Wait()

	// print the slot
	set := make(map[uint64]interface{}, 100)
	l := list.New()
	l.PushBack(slot.root)
	for {
		front := l.Front()
		if front == nil {
			break
		}
		l.Remove(front)
		ele, _ := front.Value.(*Node)
		t.Logf("element : %s", ele)
		set[ele.index] = nil
		if ele.left != nil {
			l.PushBack(ele.left)
		}
		if ele.right != nil {
			l.PushBack(ele.right)
		}
	}

	if len(set) != 100 {
		t.Logf("FAILED")
	} else {
		t.Logf("SUCCESS")
	}
}
