package storage

type Node struct {
	Key        string
	Value      string
	LeftChild  *Node
	RightChild *Node
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
type AVLTree struct {
	Root *Node
}

// Insert
// Travel the tree, compare the key with the current node key
// if it is larger -> right branch
// if it is smaller -> left branch
// if equal, replace
// if the child branch is null -> create a new node
// if the child branch is not null -> try moving down the branch
// after insert, if the height of current right and left is off by atleast 2, rebalance
// must track the parent node
// 4 case
// C1: insert to left, no right. parent has no right -> rotate parent to right
//
//	move parent to right, promote last node to parent
//
// C2: same, but right no left, parent has no left
// Q1: how to know the height of the current level ?
func (t *AVLTree) Insert(key string, value string) bool {
	node := t.Root
	for node != nil {

	}
	if key == node.Key {
		node.Key = key
		node.Value = value
		return true
	}
	return false
}
