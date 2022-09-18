package tree

type TrieNode struct {
	childNodes map[rune]*TrieNode // child nodes
	Data       string             // the leaf node will store all the string
	End        bool               // indicates that the node is leaf or not
}

type TrieTree struct {
	root *TrieNode
}

// AddChild - add the rune to the node as child.
// Return the added child node.
func (node *TrieNode) AddChild(c rune) *TrieNode {
	if node.childNodes == nil {
		node.childNodes = make(map[rune]*TrieNode)
	}

	if targetNode, ok := node.childNodes[c]; ok {
		return targetNode
	}

	node.childNodes[c] = &TrieNode{
		childNodes: nil,
		Data:       node.Data + string(c),
		End:        false,
	}
	return node.childNodes[c]
}

// FindChild - find the target rune in TrieTree
func (node *TrieNode) FindChild(c rune, hierarchy bool) *TrieNode {
	if node.childNodes == nil {
		return nil
	}

	if trieNode, ok := node.childNodes[c]; ok {
		return trieNode
	}

	if hierarchy {
		for _, trieNode := range node.childNodes {
			if targetNode := trieNode.FindChild(c, true); targetNode != nil {
				return targetNode
			}
		}
	}
	return nil
}
