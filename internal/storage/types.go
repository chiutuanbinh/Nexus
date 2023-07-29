package storage

import (
	"log"
	"strings"
)

type Node struct {
	Key           string
	Value         string
	Left          *Node
	Right         *Node
	Parent        *Node
	BalanceFactor int
}

// If there are 2 equal key, replace the value, we do not allow duplicate key
type AVLTree struct {
	Root *Node
}

func (t *AVLTree) Insert(key string, value string) error {
	if t.Root == nil {
		t.Root = &Node{
			Key:           key,
			Value:         value,
			BalanceFactor: 0,
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
			node.BalanceFactor -= 1
			if node.Left != nil {
				node = node.Left
			} else {
				node.Left = &Node{Key: key, Value: value, Parent: node}
				t.rebalance(node.Left)
				return nil
			}

		} else {
			node.BalanceFactor += 1
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

func (t *AVLTree) ToPreorderString() string {
	var stack = make([]*Node, 0)

	var sb strings.Builder

	stack = append(stack, t.Root)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top != nil {
			sb.WriteString("[" + top.Key + ":" + top.Value + "],")
			stack = append(stack, top.Right)
			stack = append(stack, top.Left)

		}
	}
	return sb.String()
}

func (n *Node) ToString() string {
	if n == nil {
		return "nil"
	}
	var sb strings.Builder
	sb.WriteString(n.Key + ":" + n.Value + " {" + n.Left.ToString() + "} {" + n.Right.ToString() + "} ")
	return sb.String()
}

func (n *Node) height() int {
	if n == nil {
		return 0
	}
	if n.Left.height() > n.Right.height() {
		return n.Left.height() + 1
	}
	return n.Right.height() + 1
}

func (t *AVLTree) Height() int {
	return t.Root.height()
}

func (n *Node) resetBalanceFactor() {
	if n == nil {
		return
	}
	n.BalanceFactor = n.Right.height() - n.Left.height()
}

func (t *AVLTree) rebalance(node *Node) {
	rn := node
	var child *Node = nil
	var grandChild *Node = nil
	for rn != nil && rn.BalanceFactor > -2 && rn.BalanceFactor < 2 {
		grandChild = child
		child = rn
		rn = rn.Parent

	}

	if rn == nil { //we reach root and tree is balance
		return
	}

	if rn.BalanceFactor == 2 { //skew right
		//child should be rn right
		if grandChild == child.Right {
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: child.Left, Parent: rn}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}
			rn.Left.resetBalanceFactor()
			rn.Key = child.Key
			rn.Value = child.Value
			rn.Right = child.Right
			rn.Right.Parent = rn

			child.Left = nil
			rn.resetBalanceFactor()

		} else {
			rn.Left = &Node{Key: rn.Key, Value: rn.Value, Left: rn.Left, Right: grandChild.Left, Parent: rn}
			if rn.Left.Left != nil {
				rn.Left.Left.Parent = rn.Left
			}
			rn.Left.resetBalanceFactor()
			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.Left = grandChild.Right
			child.Left.Parent = child
			child.resetBalanceFactor()
			grandChild.Left = nil
			grandChild.Right = nil
			rn.resetBalanceFactor()

		}
	}

	if rn.BalanceFactor == -2 { //skew left
		//child should be rn Left
		if grandChild == child.Left {
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: child.Right, Right: rn.Right, Parent: rn}
			if rn.Right.Right != nil {
				rn.Right.Right.Parent = rn.Right
			}

			rn.Right.resetBalanceFactor()
			rn.Key = child.Key
			rn.Value = child.Value
			rn.Left = child.Left
			rn.Left.Parent = rn

			child.Right = nil
			rn.resetBalanceFactor()

		} else {
			rn.Right = &Node{Key: rn.Key, Value: rn.Value, Left: grandChild.Right, Right: rn.Right, Parent: rn}
			if rn.Right.Right != nil {
				rn.Right.Right.Parent = rn.Right
			}
			rn.Right.resetBalanceFactor()
			rn.Value = grandChild.Value
			rn.Key = grandChild.Key
			child.Right = grandChild.Left
			child.Right.Parent = child
			child.resetBalanceFactor()
			grandChild.Left = nil
			grandChild.Right = nil
			rn.resetBalanceFactor()

		}
	}
	for rn != nil {
		rn.resetBalanceFactor()
		rn = rn.Parent

	}
}
