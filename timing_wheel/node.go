package timing_wheel

import "fmt"

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
