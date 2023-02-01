package timing_wheel

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var L0ArrangeError = errors.New("current slot already in L1 bucket. can't rearrange")
var ArrangeEOF = errors.New("arrange node to slot occurs EOF")

// Slot - AVL tree
type Slot struct {
	root    *Node
	current *Node   // record the new node should be added after the current node.
	bucket  *Bucket // record the slot that in the bucket.
}

// Append - add node to the end of Slot without lock.
func (slot *Slot) Append(node *Node) {
	// if the root is nil, then record the node as root.
	if slot.root == nil {
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&slot.root)), nil, unsafe.Pointer(node)) {
			atomic.AddInt32(&node.depth, 1)
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&slot.current)), unsafe.Pointer(node))
			return
		}
	}
	for {
		cur := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&slot.current))))
		if cur == nil {
			continue
		}
		if cur != slot.root && cur.parent.left == cur {
			if slot.leftAppend(cur, node) {
				break
			}
		} else {
			if slot.rightAppend(cur, node) {
				break
			}
		}
	}

}

// leftAppend - current is the left leaf of the parent.
func (slot *Slot) leftAppend(cur, node *Node) bool {
	if atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&cur.next)), nil, unsafe.Pointer(node),
	) {
		cur.parent.right = node
		node.parent = cur.parent
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&slot.current)), unsafe.Pointer(node))
		return true
	}
	return false
}

// rightAppend - current is the right leaf of the parent.
func (slot *Slot) rightAppend(cur, node *Node) bool {
	// check whether the current node is the end of the layer of the tree.
	if cur == slot.root || cur.parent.next == nil {
		// the node added to the first node of the layer of the tree
		left := slot.root
		right := slot.root
		for {
			if right == cur {
				break
			}
			left = left.left
			right = right.right
		}
		node.parent = left
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&left.left)), nil, unsafe.Pointer(node),
		) {
			atomic.AddInt32(&slot.root.depth, 1)
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&slot.current)), unsafe.Pointer(node))
			return true
		}
	} else {
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&cur.next)), nil, unsafe.Pointer(node),
		) {
			cur.parent.next.left = node
			node.parent = cur.parent.next
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&slot.current)), unsafe.Pointer(node))
			return true
		}
	}
	return false
}

func (slot *Slot) ArrangeTo(l0 *Bucket) error {
	if slot.bucket == l0 {
		return L0ArrangeError
	}
	//maxTotalNodes := int(math.Pow(2.0, float64(slot.root.depth))) - 1
	// evaluate the number of goroutines that should call to locate the node.
	goroutines := int(slot.root.depth - 8)
	if goroutines > runtime.NumCPU() {
		goroutines = runtime.NumCPU()
	}
	// relocate the nodes in threads.

	return nil
}

func (slot *Slot) arrangeNodeTo(node *Node, l0 *Bucket) error {
	lvl, _, _ := location(node.index)
	round := 256
	for i := 1; i < lvl; i++ {
		round *= 64
	}
	lvl, lvl_idx, idx := location(node.index - uint64(round))
	if lvl == 0 {
		l0.slots[idx].Append(node)
		return nil
	}
	bucket := l0
	for {
		if bucket == nil || slot.bucket == bucket {
			return ArrangeEOF
		}
		if lvl == 0 {
			bucket.slots[lvl_idx].Append(node)
			return nil
		}
		lvl--
		bucket = bucket.next
	}
	return nil
}
