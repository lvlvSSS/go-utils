package timing_wheel

import (
	"sync/atomic"
	"unsafe"
)

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
