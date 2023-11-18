package storage

import (
	"math"
	"strings"

	"github.com/rs/zerolog/log"
)

type Node struct {
	Key    string
	Value  string
	Left   *Node
	Right  *Node
	Parent *Node
}

type Memtable interface {
	Insert(key string, value string) error
	Delete(key string) (bool, error)
	Find(key string) (string, bool)
	Clear() error
	List() []Tuple
	Size() int
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
type AVLTree struct {
	Root *Node
	size int
}

func (t *AVLTree) Insert(key string, value string) error {
	t.size += len(key) + len(value)
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
		if key == node.Key {
			node.Key = key
			node.Value = value
			return nil
		}
		if key < node.Key {
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

func (t *AVLTree) Delete(key string) (bool, error) {
	node := t.Root
	for {
		if node == nil {
			return false, nil
		}
		if key == node.Key {
			t.size -= len(node.Key) + len(node.Value)
			if node == t.Root {
				if node.Left == nil {
					t.Root = node.Right
					t.Root.Parent = nil
					return true, nil
				}
				if node.Right == nil {
					t.Root = node.Left
					t.Root.Parent = nil
					return true, nil
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
					return true, nil
				}
				if node.Right == nil {
					if node == parent.Right {
						parent.Right = node.Left
					} else {
						parent.Left = node.Left
					}
					t.rebalance(parent)
					return true, nil
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
			return true, nil
		}
		if key < node.Key {
			node = node.Left
		} else {
			node = node.Right
		}
	}
}

func (t *AVLTree) Find(key string) (string, bool) {
	node := t.Root
	for {
		if node == nil {
			return "", false
		}
		if key < node.Key {
			node = node.Left
		} else if key > node.Key {
			node = node.Right
		} else {
			return node.Value, true
		}
	}
}

func (t *AVLTree) Clear() error {
	log.Debug().Msg("AVL clear")
	t.Root = nil
	t.size = 0
	return nil
}

func (t *AVLTree) Size() int {
	return t.size
}

type Tuple struct {
	Key   string
	Value string
}

func (t *AVLTree) List() []Tuple {
	return t.Inorder()
}

func (t *AVLTree) Inorder() []Tuple {
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
	sb.WriteString(n.Key)
	if n.Parent != nil {
		sb.WriteString(":" + "{" + n.Parent.Key + "}")
	} else {
		sb.WriteString(":ROOT")
	}
	sb.WriteString(" {" + n.Left.ToString() + "} {" + n.Right.ToString() + "}")
	return sb.String()
}

func (t *AVLTree) Height() int {
	return t.Root.Height()
}

func (n *Node) Height() int {
	if n == nil {
		return 0
	}
	return 1 + int(math.Max(float64(n.Left.Height()), float64(n.Right.Height())))
}

func (t *AVLTree) rebalance(node *Node) {
	rn := node
	var child *Node = nil
	var grandChild *Node = nil
	var skewLeft, skewRight bool
	for rn != nil {
		skewLeft = rn.skewLeft()
		skewRight = rn.skewRight()
		if skewLeft || skewRight {
			break
		}
		grandChild = child
		child = rn
		rn = rn.Parent

	}

	if rn == nil { //we reach root and tree is balance
		return
	}

	if skewRight {
		//child should be rn right
		if grandChild == child.Right {
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: child.Left, Parent: rn}
			if child.Left != nil {
				child.Left.Parent = rn.Left
			}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}

			rn.Key = child.Key
			rn.Value = child.Value
			rn.setRight((child.Right))

			child.Left = nil

		} else {
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: grandChild.Left, Parent: rn}
			if grandChild.Left != nil {
				grandChild.Left.Parent = rn.Left
			}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}

			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.setLeft(grandChild.Right)

			grandChild.Left = nil
			grandChild.Right = nil

		}
	} else if skewLeft {
		//child should be rn Left
		if grandChild == child.Left {
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: child.Right, Right: rn.Right, Parent: rn}
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

			child.Right = nil

		} else {
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: grandChild.Right, Right: rn.Right, Parent: rn}
			if grandChild.Right != nil {
				grandChild.Right.Parent = rn.Right
			}
			if rn.Right.Right != nil {
				rn.Right.Right.Parent = rn.Right
			}

			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.setRight(grandChild.Left)

			grandChild.Left = nil
			grandChild.Right = nil

		}
	}
}

func (n *Node) skewRight() bool {
	return n.Right.Height()-n.Left.Height() >= 2
}

func (n *Node) skewLeft() bool {
	return n.Left.Height()-n.Right.Height() >= 2
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
