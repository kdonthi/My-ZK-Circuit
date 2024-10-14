package succinct

import (
	"fmt"
	"math"
)

type HintOp int

const (
	Subtraction HintOp = iota
	Addition
	Multiplication
	Division
	Square
	Sqrt
)

// c = a + 3
// b = c - 2 => b + 2 = c
// BUT Hint: b = c - 2

type HintBuilder struct {
	n                     *NodeBuilder
	first, second, output *Node
}

func NewHint() *HintBuilder {
	return &HintBuilder{
		n: NewNodeBuilder(),
	}
}

type HintNode struct {
	typ          NodeType
	val          float64
	op           HintOp
	children     []*HintNode
	dependencies map[int]*Node
	id           int
}

func mergeDeps(m1, m2 map[int]*Node) map[int]*Node {
	m := map[int]*Node{}

	for id, elem := range m1 {
		m[id] = elem
	}
	for id, elem := range m2 {
		m[id] = elem
	}

	return m
}

func (h *HintBuilder) Val(n *Node) *HintNode {
	return &HintNode{
		typ: Number,
		val: n.val,
		dependencies: map[int]*Node{
			n.id: n,
		},
		children: make([]*HintNode, 2),
	}
}

func (h *HintBuilder) Const(n float64) *HintNode {
	return &HintNode{
		typ:          Number,
		val:          n,
		dependencies: make(map[int]*Node),
		children:     make([]*HintNode, 2),
	}
}

func (h *HintBuilder) Add(n1, n2 *HintNode) *HintNode {
	return &HintNode{
		typ:          Operation,
		op:           Addition,
		children:     []*HintNode{n1, n2},
		dependencies: mergeDeps(n1.dependencies, n2.dependencies),
	}
}

func (h *HintBuilder) Sub(n1, n2 *HintNode) *HintNode {
	return &HintNode{
		typ:          Operation,
		op:           Subtraction,
		children:     []*HintNode{n1, n2},
		dependencies: mergeDeps(n1.dependencies, n2.dependencies),
	}
}

func (h *HintBuilder) Mul(n1, n2 *HintNode) *HintNode {
	return &HintNode{
		typ:          Operation,
		op:           Multiplication,
		children:     []*HintNode{n1, n2},
		dependencies: mergeDeps(n1.dependencies, n2.dependencies),
	}
}

func (h *HintBuilder) Div(n1, n2 *HintNode) *HintNode {
	return &HintNode{
		typ:          Operation,
		op:           Division,
		children:     []*HintNode{n1, n2},
		dependencies: mergeDeps(n1.dependencies, n2.dependencies),
	}
}

func (h *HintBuilder) Square(n1 *HintNode) *HintNode {
	return &HintNode{
		typ:          Operation,
		op:           Square,
		children:     []*HintNode{n1, nil},
		dependencies: n1.dependencies,
	}
}

func (h *HintBuilder) Sqrt(n1 *HintNode) *HintNode {
	return &HintNode{
		typ:          Operation,
		op:           Sqrt,
		children:     []*HintNode{n1, nil},
		dependencies: n1.dependencies,
	}
}

func (h *HintNode) Build(n *Node) *Hint {
	return &Hint{
		output:   n, // TODO do we need this?
		equation: h,
	}
}

// hint should have the node ids that it uses
type Hint struct {
	output   *Node
	equation *HintNode
}

// solve with the map you were given
func (h *Hint) Solve(m map[int]*Node) bool {
	for id, v := range h.equation.dependencies { // TODO I don't knw if it makes sense to do MaybeInts?
		if !v.valFilled {
			panic(fmt.Sprintf("not all variables solved, e.g. %v", id))
		}
	}

	var s []*HintNode
	var s2 []*HintNode

	s = append(s, h.equation)
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
			v := m[lastElem.id]
			lastElem.val = v.val
		} else if lastElem.typ == Operation {
			first := lastElem.children[0].val
			second := lastElem.children[1].val

			switch lastElem.op {
			case Addition:
				lastElem.val = first + second
			case Multiplication:
				lastElem.val = first * second
			case Subtraction:
				lastElem.val = first - second
			case Division:
				lastElem.val = first / second
			case Square:
				lastElem.val = first * first
			case Sqrt:
				lastElem.val = math.Sqrt(first)
			}
		}
	}

	return true
}

// hint should have a "Solve" function that allows it to have a true or false!
// maybe to have a list that keeps track of the nodes we don't know the value of, and if the size doesn't change across iterations, we say it's an error?!!!
func (h *Hint) Build(output *Node) *Hint {
	h.output = output
	return h
}
