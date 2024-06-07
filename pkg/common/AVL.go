package common

import (
	"bytes"
	"math"
	"strings"

	"github.com/rs/zerolog/log"
)

type Node struct {
	Key           []byte
	Value         []byte
	Left          *Node
	Right         *Node
	Parent        *Node
	BalanceFactor int8
}

func CreateAVLTree() BstTree {
	return &avlTreeImpl{}

}

type avlTreeImpl struct {
	Root      *Node
	size      int
	nodeCount int
}

// LowerBound return the node with smallest key that equal or larger than key
func (t *avlTreeImpl) LowerBound(key []byte) ([]byte, []byte, error) {
	node := t.Root
	var k []byte = nil
	var v []byte = nil
	for {
		if node == nil {
			return k, v, nil
		}
		if bytes.Equal(key, node.Key) {
			return node.Key, node.Value, nil
		} else if bytes.Compare(key, node.Key) == 1 {
			node = node.Right
		} else {
			if k == nil {
				k = node.Key
				v = node.Value
			} else {
				if bytes.Compare(k, node.Key) == 1 {
					k = node.Key
					v = node.Value
				}
			}
			node = node.Left
		}
	}

}

// UpperBound return the node with smallest key that larger than key
func (t *avlTreeImpl) UpperBound(key []byte) ([]byte, []byte, error) {
	node := t.Root
	var k []byte = nil
	var v []byte = nil
	for {
		if node == nil || bytes.Equal(key, node.Key) {
			return k, v, nil
		} else if bytes.Compare(key, node.Key) == 1 {
			node = node.Right
		} else {
			if k == nil {
				k = node.Key
				v = node.Value
			} else {
				if bytes.Compare(k, node.Key) == 1 {
					k = node.Key
					v = node.Value
				}
			}
			node = node.Left
		}
	}
}

func (t *avlTreeImpl) Insert(key []byte, value []byte) error {
	// log.Info().Msgf("Tree height %v", t.Height())
	t.size += len(key) + len(value)
	t.nodeCount += 1
	if t.Root == nil {
		t.Root = &Node{
			Key:   key,
			Value: value,
		}
		return nil
	}
	node := t.Root

	for {
		if node == nil {
			return nil
		}
		if bytes.Equal(key, node.Key) {
			node.Key = key
			node.Value = value
			return nil
		}
		if bytes.Compare(key, node.Key) == -1 {
			if node.Left != nil {
				node = node.Left
			} else {
				node.Left = &Node{Key: key, Value: value, Parent: node}
				t.rebalance(node.Left)
				return nil
			}

		} else {
			if node.Right != nil {
				node = node.Right
			} else {
				node.Right = &Node{Key: key, Value: value, Parent: node}
				t.rebalance(node.Right)
				return nil
			}
		}

	}
}

func (t *avlTreeImpl) Delete(key []byte) error {
	node := t.Root
	for {
		if node == nil {
			return NotFoundError{Value: string(key)}
		}
		if bytes.Equal(key, node.Key) {
			t.size -= len(node.Key) + len(node.Value)
			t.nodeCount -= 1
			if node == t.Root {
				if node.Left == nil {
					t.Root = node.Right
					t.Root.Parent = nil
					return nil
				}
				if node.Right == nil {
					t.Root = node.Left
					t.Root.Parent = nil
					return nil
				}

			} else {
				parent := node.Parent
				if node.Left == nil {
					if node == parent.Left {
						parent.setLeft(node.Right)
					} else {
						parent.setRight(node.Right)
					}
					t.rebalance(parent)
					return nil
				}
				if node.Right == nil {
					if node == parent.Right {
						parent.Right = node.Left
					} else {
						parent.Left = node.Left
					}
					t.rebalance(parent)
					return nil
				}
			}

			candidate := node.Right.leftMostDescendant()

			node.Value = candidate.Value
			node.Key = candidate.Key

			cp := candidate.Parent
			if cp == node {
				node.setRight(candidate.Right)
			} else {
				cp.setLeft(candidate.Right)

			}
			t.rebalance(cp)
			return nil
		}
		if bytes.Compare(key, node.Key) == -1 {
			node = node.Left
		} else {
			node = node.Right
		}
	}
}

func (t *avlTreeImpl) Find(key []byte) ([]byte, error) {
	node := t.Root
	for {
		if node == nil {
			return nil, NotFoundError{Value: string(key)}
		}
		if bytes.Compare(key, node.Key) == -1 {
			node = node.Left
		} else if bytes.Compare(key, node.Key) == 1 {
			node = node.Right
		} else {
			return node.Value, nil
		}
	}
}

func (t *avlTreeImpl) Clear() error {
	log.Debug().Msg("AVL clear")
	t.Root = nil
	t.size = 0
	return nil
}

func (t *avlTreeImpl) Size() int {
	return t.size
}

func (t *avlTreeImpl) NodeCount() int {
	return t.nodeCount
}

func (t *avlTreeImpl) List() []Tuple {
	return t.Inorder()
}

func (t *avlTreeImpl) Inorder() []Tuple {
	return t.Root.Inorder()
}

func (n *Node) Inorder() []Tuple {
	if n == nil {
		return make([]Tuple, 0)
	}
	res := make([]Tuple, 0)
	res = append(res, n.Left.Inorder()...)
	res = append(res, Tuple{Key: n.Key, Value: n.Value})
	res = append(res, n.Right.Inorder()...)
	return res
}

func (n *Node) ToString() string {
	if n == nil {
		return "nil"
	}
	var sb strings.Builder
	sb.WriteString(string(n.Key))
	if n.Parent != nil {
		sb.WriteString(":" + "{" + string(n.Parent.Key) + "}")
	} else {
		sb.WriteString(":ROOT")
	}
	sb.WriteString(" {" + n.Left.ToString() + "} {" + n.Right.ToString() + "}")
	return sb.String()
}

func (t *avlTreeImpl) Height() int {
	return t.Root.Height()
}

func (n *Node) Height() int {
	if n == nil {
		return 0
	}
	return 1 + int(math.Max(float64(n.Left.Height()), float64(n.Right.Height())))
}

func (t *avlTreeImpl) rebalance(node *Node) {

	rn := node
	var child *Node = nil
	var grandChild *Node = nil
	var change = true
	for {
		if rn.Parent == nil {
			break
		}

		if rn == rn.Parent.Right {
			rn.Parent.BalanceFactor += 1
			if rn.Parent.BalanceFactor == 0 {
				change = false
			}
		} else {
			rn.Parent.BalanceFactor -= 1
			if rn.Parent.BalanceFactor == 0 {
				change = false
			}
		}
		grandChild = child
		child = rn
		rn = rn.Parent
		if rn == nil {
			return
		}
		if rn.skewLeft() || rn.skewRight() || !change {
			break
		}
	}

	if rn.skewRight() {
		//child should be rn right
		if grandChild == child.Right { //right right

			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: child.Left, Parent: rn, BalanceFactor: 0}
			if child.Left != nil {
				child.Left.Parent = rn.Left
			}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}

			rn.Key = child.Key
			rn.Value = child.Value
			rn.setRight((child.Right))
			rn.BalanceFactor = 0
			child.Left = nil

		} else { //right left
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: grandChild.Left, Parent: rn, BalanceFactor: 0}
			if grandChild.Left != nil {
				grandChild.Left.Parent = rn.Left
			}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}

			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.setLeft(grandChild.Right)
			rn.BalanceFactor = 0
			grandChild.Left = nil
			grandChild.Right = nil

		}
	} else if rn.skewLeft() {
		//child should be rn Left
		if grandChild == child.Left { //left right
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: child.Right, Right: rn.Right, Parent: rn, BalanceFactor: 0}
			if child.Right != nil {
				child.Right.Parent = rn.Right
			}
			if rn.Right.Right != nil {
				rn.Right.Right.Parent = rn.Right
			}
			rn.Right.setRight(rn.Right.Right)

			rn.Key = child.Key
			rn.Value = child.Value
			rn.setLeft(child.Left)
			rn.BalanceFactor = 0
			child.Right = nil

		} else { //right left
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: grandChild.Right, Right: rn.Right, Parent: rn, BalanceFactor: 0}
			if grandChild.Right != nil {
				grandChild.Right.Parent = rn.Right
			}
			if rn.Right.Right != nil {
				rn.Right.Right.Parent = rn.Right
			}

			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.setRight(grandChild.Left)
			rn.BalanceFactor = 0
			grandChild.Left = nil
			grandChild.Right = nil

		}
	}
}

func (n *Node) skewRight() bool {
	return n.BalanceFactor == 2
}

func (n *Node) skewLeft() bool {
	return n.BalanceFactor == -2
}

func (n *Node) leftMostDescendant() *Node {
	res := n
	for res.Left != nil {
		res = res.Left
	}
	return res
}

func (n *Node) setLeft(child *Node) bool {
	n.Left = child
	if child != nil {
		child.Parent = n
	}
	return true
}

func (n *Node) setRight(child *Node) bool {
	n.Right = child
	if child != nil {
		child.Parent = n
	}
	return true
}

func (n *Node) Size() int {
	return len(n.Key) + len(n.Value)
}
