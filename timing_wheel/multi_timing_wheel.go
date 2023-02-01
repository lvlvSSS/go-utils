package timing_wheel

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var Unmodified = errors.New("state in timing wheel can't be modified now")

type state uint32

const (
	unstarted state = 1 << iota
	running
	finished
)

type TimingWheel struct {
	bucket *Bucket

	jiffy time.Duration
	next  *TimingWheel

	state state
	lock  sync.Mutex
}

func (wheel *TimingWheel) DefineJiffy(dur time.Duration) error {
	if state(atomic.LoadUint32((*uint32)(&wheel.state))) == unstarted {
		wheel.lock.Lock()
		defer wheel.lock.Unlock()
		if state(atomic.LoadUint32((*uint32)(&wheel.state))) == unstarted {
			wheel.jiffy = dur
		}
		return nil
	}
	return Unmodified
}

func (wheel *TimingWheel) Run() error {
	for {

		time.Sleep(wheel.jiffy / 4)
	}
	return nil
}

type Node struct {
	index int
	do    func()

	// depth, parent, left, right, next - all used to iterate the tree that built by Nodes.
	depth       int32 // only the root node record the value.
	parent      *Node
	left, right *Node
	next        *Node
}

func (node *Node) String() string {
	return fmt.Sprintf("Node{index[%d], depth[%d], next[%v]}", node.index, node.depth, node.next)
}

type Bucket struct {
	slots []Slot
	next  *Bucket
}

// Slot - AVL tree
type Slot struct {
	root    *Node
	current *Node // record the new node should be added after the current node.
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

// location - like linux's timing wheel, l1 has 256 slots, l2~ln has 64 slots.
func location(duration uint64) (lvl, lvl_idx, idx int) {
	const L1 = 0xff
	const Ln = 0x3f
	idx = int(duration & L1)
	lvl = 0
	for cur := duration >> 8; cur&Ln > 0 || cur > Ln; cur = cur >> 6 {
		lvl_idx = int(cur&Ln) - 1
		lvl++
	}
	/*	cur := jiffy >> 8
		for {
			if cur&Ln <= 0 && cur < Ln {
				return
			}
			lvl++
			cur = cur >> 6
		}*/
	return
}
