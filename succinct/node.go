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

func (n *Node) GetVal() float64 {
	return n.val
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
	hints      map[int]*Hint
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

	nb.Solve(head)
	for _, eq := range nb.equalities {
		nb.Solve(eq.n1)

		val1 := float64(0)
		if eq.n1.valFilled {
			if v, ok := nb.m[eq.n1.id]; ok {
				val1 = v.i
			} else {
				return false
			}
		}

		nb.Solve(eq.n2)

		val2 := float64(0)
		if eq.n2.valFilled {
			if v, ok := nb.m[eq.n2.id]; ok {
				val2 = v.i
			} else {
				return false
			}
		}

		if val1 != val2 {
			return false
		}
	}

	return true
}

func (nb *NodeBuilder) Solve(head *Node) {
	var s []*Node
	var s2 []*Node

	s = append(s, head)
	for len(s) != 0 {
		lastElem := s[len(s)-1]
		s = s[:len(s)-1]
		if lastElem.valFilled {
			continue
		}

		s2 = append(s2, lastElem)

		if len(lastElem.children) != 0 {
			left := lastElem.children[0]
			right := lastElem.children[1]

			s = append(s, right)
			s = append(s, left)
		}
	}

	refill := []*Node{}
	for {
		for len(s2) != 0 {
			lastElem := s2[len(s2)-1]
			s2 = s2[:len(s2)-1]

			if lastElem.typ == Variable {
				v := nb.m[lastElem.id]
				lastElem.val = v.i
				lastElem.valFilled = true
			} else if lastElem.typ == Operation {
				node1 := lastElem.children[0]
				node2 := lastElem.children[1]

				if !node1.valFilled || !node2.valFilled {
					refill = append(refill, lastElem)
					continue
				}

				first := node1.val
				second := node2.val

				switch lastElem.op {
				case Add:
					lastElem.val = first + second
				case Mul:
					lastElem.val = first * second
				}
			} else if lastElem.typ == Hinted {
				if hint, ok := nb.hints[lastElem.id]; ok {
					if val, ok := hint.Solve(); ok {
						lastElem.val = val
						lastElem.valFilled = true
					} else {
						refill = append(refill, lastElem)
					}
				} else {
					panic("hint needed")
				}
			}
		}

		if len(refill) == 0 {
			break
		}

		copy(s2, refill)
		refill = []*Node{}
	}
}

func (nb *NodeBuilder) Hint() *HintBuilder {
	return nb.h
}
