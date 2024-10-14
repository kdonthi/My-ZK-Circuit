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
	Hinted
)

type Op int

const (
	Add Op = iota
	Mul
)

type Node struct {
	typ       NodeType
	val       float64
	valFilled bool
	op        Op
	children  []*Node
	id        int
}

func (n *Node) GetId() int {
	return n.id
}

type MaybeInt struct {
	i float64
	b bool
}

type Equalities struct {
	n1, n2 *Node
}

type NodeBuilder struct {
	m          map[int]MaybeInt
	equalities []Equalities
	hint       []Hint
	h          *HintBuilder
}

func NewNodeBuilder() *NodeBuilder {
	return &NodeBuilder{
		m:          map[int]MaybeInt{},
		equalities: make([]Equalities, 0),
		h:          NewHintBuilder(),
	}
}

func (nb *NodeBuilder) Var() *Node {
	id := c.GetNext()
	nb.m[id] = MaybeInt{
		b: false,
	}

	return &Node{
		typ:       Variable,
		id:        id,
		valFilled: true,
	}
}

func (nb *NodeBuilder) Num(n float64) *Node {
	return &Node{
		typ:       Number,
		val:       n,
		valFilled: true,
	}
}

func (nb *NodeBuilder) Add(n1, n2 *Node) *Node {
	return &Node{
		typ:       Operation,
		op:        Add,
		children:  []*Node{n1, n2},
		valFilled: false,
	}
}

func (nb *NodeBuilder) Mul(n1, n2 *Node) *Node {
	return &Node{
		typ:       Operation,
		op:        Mul,
		children:  []*Node{n1, n2},
		valFilled: false,
	}
}

func (nb *NodeBuilder) AssertEq(n1, n2 *Node) {
	nb.equalities = append(nb.equalities, Equalities{
		n1: n1,
		n2: n2,
	})
}

func (nb *NodeBuilder) Const(val float64) *Node {
	return &Node{
		typ:       Number,
		val:       val,
		valFilled: true,
	}
}

func (nb *NodeBuilder) FillNodes(m map[int]float64) {
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
	for id, v := range nb.m { // TODO I don't knw if it makes sense to do MaybeInts?
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
		val1 := eq.n1.val
		if v, ok := nb.m[eq.n1.id]; ok {
			val1 = v.i
		}

		val2 := eq.n2.val
		if v, ok := nb.m[eq.n2.id]; ok {
			val2 = v.i
		}

		if val1 != val2 {
			return false
		}
	}

	return true
}

func (nb *NodeBuilder) Hint() *HintBuilder {
	return nb.h
}
