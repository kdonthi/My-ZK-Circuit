package succinct

import (
	"awesomeProject/utils"
	"fmt"
)

var c *utils.Counter

func init() {
	c = utils.NewCounter()
}

type NodeType int

const (
	Number NodeType = iota
	Operation
	Variable
	NoType
)

type Op int

const (
	Add Op = iota
	Mul
	NoOp
)

type Node struct {
	typ      NodeType
	val      int
	op       Op
	children []*Node
	id       int
}

func (n *Node) GetId() int {
	return n.id
}

type MaybeInt struct {
	i int
	b bool
}

type Equalities struct {
	n1, n2 *Node
}

type NodeBuilder struct {
	m          map[int]MaybeInt
	equalities []Equalities
}

func NewNodeBuilder() NodeBuilder {
	return NodeBuilder{
		m:          map[int]MaybeInt{},
		equalities: make([]Equalities, 0),
	}
}

func (nb *NodeBuilder) Var() *Node {
	id := c.GetNext()
	nb.m[id] = MaybeInt{
		b: false,
	}

	return &Node{
		typ: Variable,
		id:  id,
	}
}

func (nb *NodeBuilder) Num(n int) *Node {
	return &Node{
		typ: Number,
		val: n,
	}
}

func (nb *NodeBuilder) Add(n1, n2 *Node) *Node {
	return &Node{
		typ:      Operation,
		op:       Add,
		children: []*Node{n1, n2},
	}
}

func (nb *NodeBuilder) Mul(n1, n2 *Node) *Node {
	return &Node{
		typ:      Operation,
		op:       Mul,
		children: []*Node{n1, n2},
	}
}

func (nb *NodeBuilder) AssertEq(n1, n2 *Node) {
	nb.equalities = append(nb.equalities, Equalities{
		n1: n1,
		n2: n2,
	})
}

func (nb *NodeBuilder) Const(val int) *Node {
	return &Node{
		typ: Number,
		val: val,
	}
}

func (nb *NodeBuilder) FillNodes(m map[int]int) {
	for v := range nb.m {
		if val, ok := m[v]; ok {
			nb.m[v] = MaybeInt{
				i: val,
				b: ok,
			}
		} else {
			panic(fmt.Sprintf("no value for id %v in map", v))
		}
	}
}

func (nb *NodeBuilder) Verify(head *Node) bool {
	for id, v := range nb.m {
		if !v.b {
			panic(fmt.Sprintf("not all variables filled, e.g. %v", id))
		}
	}

	var s []*Node
	var s2 []*Node

	s = append(s, head)
	for len(s) != 0 {
		lastElem := s[len(s)-1]
		s2 = append(s2, lastElem)

		s = s[:len(s)-1]

		if len(lastElem.children) != 0 {
			left := lastElem.children[0]
			right := lastElem.children[1]

			s = append(s, right)
			s = append(s, left)
		}
	}

	for len(s2) != 0 {
		lastElem := s2[len(s2)-1]
		s2 = s2[:len(s2)-1]
		if lastElem.typ == Variable {
			v := nb.m[lastElem.id]
			lastElem.val = v.i
		} else if lastElem.typ == Operation {
			first := lastElem.children[0].val
			second := lastElem.children[1].val

			switch lastElem.op {
			case Add:
				lastElem.val = first + second
			case Mul:
				lastElem.val = first * second
			}
		}
	}

	for _, eq := range nb.equalities {
		val1 := eq.n1.id
		if eq.n1.val != eq.n2.val {
			return false
		}
	}

	return true
}
