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

// HintBuilder supports the operations to create a hint (use "Build" on created Hint to create a Node)
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
	wrappedNode  *Node // underlying node for hinted values that depend on a real node
	op           HintOp
	children     []*HintNode
	dependencies map[int]*Node // nodes this value depends on to be able to solve
	solved       bool
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
	if n2.val == 0 {
		panic("cannot divide by 0")
	}

	return &HintNode{
		typ:          Operation,
		op:           Division,
		children:     []*HintNode{n1, n2},
		dependencies: mergeDeps(n1.dependencies, n2.dependencies),
	}
}

func (h *HintBuilder) Sqrt(n1 *HintNode) *HintNode {
	if n1.val < 0 {
		panic("cannot take the square root of a negative number (at least in this circuit :))")
	}
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

func (h *HintNode) solve() (float64, bool) {
	for _, v := range h.dependencies {
		if !v.valFilled && v.typ != Hinted {
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
		} else if lastElem.typ == Hinted {
			if !lastElem.solved {
				if val, ok := lastElem.solve(); ok {
					lastElem.val = val
					lastElem.solved = true
				} else {
					return 0, false
				}
			}
		}
	}

	return h.val, true
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
