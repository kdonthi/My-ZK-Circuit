package succinct

import (
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
	hints map[int]*HintNode
}

func NewHintBuilder() *HintBuilder {
	return &HintBuilder{
		hints: map[int]*HintNode{},
	}
}

type HintNode struct {
	typ          NodeType
	val          float64
	wrappedNode  *Node
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
		typ:         Variable,
		wrappedNode: n,
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
		children:     make([]*HintNode, 0),
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

func (h *HintBuilder) Build(n *HintNode) *Node {
	id := c.GetNext()
	h.hints[id] = n

	return &Node{
		typ:       Hinted,
		valFilled: false,
		id:        id,
	}
}

// hint should have the node ids that it uses
type Hint struct {
	equation *HintNode
}

// solve with the map you were given
func (h *HintNode) Solve() (float64, bool) {
	for _, v := range h.dependencies { // TODO I don't knw if it makes sense to do MaybeInts?
		if !v.valFilled {
			return 0, false
		}
	}

	var s []*HintNode
	var s2 []*HintNode

	s = append(s, h)
	for len(s) != 0 {
		lastElem := s[len(s)-1]
		s2 = append(s2, lastElem)

		s = s[:len(s)-1]

		if len(lastElem.children) != 0 {
			left := lastElem.children[0]
			right := lastElem.children[1]

			if left != nil {
				s = append(s, left)
			}
			if right != nil {
				s = append(s, right)
			}
		}
	}

	for len(s2) != 0 {
		lastElem := s2[len(s2)-1]
		s2 = s2[:len(s2)-1]
		if lastElem.typ == Variable {
			if lastElem.wrappedNode.valFilled {
				lastElem.val = lastElem.wrappedNode.val
			}
		} else if lastElem.typ == Operation {
			var first, second *HintNode
			first = lastElem.children[0]
			second = lastElem.children[1]

			switch lastElem.op {
			case Addition:
				lastElem.val = first.val + second.val
			case Multiplication:
				lastElem.val = first.val * second.val
			case Subtraction:
				lastElem.val = first.val - second.val
			case Division:
				lastElem.val = first.val / second.val
			case Square:
				lastElem.val = first.val * first.val
			case Sqrt:
				lastElem.val = math.Sqrt(first.val)
			}
		}
	}

	return h.val, true
}

// hint should have a "Solve" function that allows it to have a true or false!
// maybe to have a list that keeps track of the nodes we don't know the value of, and if the size doesn't change across iterations, we say it's an error?!!!
//func (h *Hint) Equals(n *Node) *Node {
//	h.output = n
//	return &Node{
//		typ:       Hinted,
//		valFilled: false,
//	}
//}
