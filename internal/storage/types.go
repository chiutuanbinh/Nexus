package storage

import (
	"log"
	"strings"
)

type Node struct {
	Key    string
	Value  string
	Left   *Node
	Right  *Node
	Parent *Node
	height int
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
type AVLTree struct {
	Root *Node
}

func (t *AVLTree) Insert(key string, value string) error {
	if t.Root == nil {
		t.Root = &Node{
			Key:    key,
			Value:  value,
			height: 1,
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
		log.Printf("GO TO %v with value %v", node, key)
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

func (t *AVLTree) Delete(key string) error {
	node := t.Root
	for {
		if node == nil {
			return nil
		}
		if key == node.Key {
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
			log.Printf("Candidate left most of right %v cp %v node %v", candidate, cp, node)
			if cp == node {
				node.setRight(candidate.Right)
			} else {
				cp.setLeft(candidate.Right)

			}
			t.rebalance(cp)
			return nil
		}
		if key < node.Key {
			node = node.Left
		} else {
			node = node.Right
		}
	}
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
	return n.height
}

func (n *Node) resetHeight() {
	if n == nil {
		return
	}
	left := n.Left.Height()
	right := n.Right.Height()
	if left > right {
		n.height = left + 1
	} else {
		n.height = right + 1
	}
}

func (t *AVLTree) rebalance(node *Node) {
	rn := node
	var child *Node = nil
	var grandChild *Node = nil
	for rn != nil {
		oldHeight := rn.Height()
		rn.resetHeight()
		if oldHeight == rn.Height() || rn.skewLeft() || rn.skewRight() {
			break
		}
		grandChild = child
		child = rn
		rn = rn.Parent

	}
	log.Printf("rebalance %v %v %v", rn, child, grandChild)

	if rn == nil { //we reach root and tree is balance
		return
	}

	if rn.skewRight() {
		//child should be rn right
		if grandChild == child.Right {
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: child.Left, Parent: rn}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}
			rn.Left.resetHeight()
			rn.Key = child.Key
			rn.Value = child.Value
			rn.setRight((child.Right))

			child.Left = nil
			rn.resetHeight()

		} else {
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: grandChild.Left, Parent: rn}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}
			rn.Left.resetHeight()
			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.setLeft(grandChild.Right)
			child.resetHeight()
			grandChild.Left = nil
			grandChild.Right = nil
			rn.resetHeight()

		}
	}
	log.Printf("L %v and R %v", rn.Left.Height(), rn.Right.Height())
	if rn.skewLeft() {
		//child should be rn Left
		if grandChild == child.Left {
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: child.Right, Right: rn.Right, Parent: rn}
			rn.Right.setRight(rn.Right.Right)

			rn.Right.resetHeight()
			rn.Key = child.Key
			rn.Value = child.Value
			rn.setLeft(child.Left)

			child.Right = nil
			rn.resetHeight()

		} else {
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: grandChild.Right, Right: rn.Right, Parent: rn}
			if rn.Right.Right != nil {
				rn.Right.Right.Parent = rn.Right
			}
			rn.Right.resetHeight()
			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.setRight(grandChild.Left)
			child.resetHeight()
			grandChild.Left = nil
			grandChild.Right = nil
			rn.resetHeight()

		}
	}
	for rn != nil {
		rn.resetHeight()
		rn = rn.Parent

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

func (n *Node) rightMostDescendant() *Node {
	res := n
	for res.Right != nil {
		res = res.Right
	}
	return res
}

func (n *Node) unlinkFromParent() bool {
	if n == nil {
		return false
	}

	if n == n.Parent.Left {
		n.Parent.Left = nil
	} else {
		n.Parent.Right = nil
	}
	return true
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
