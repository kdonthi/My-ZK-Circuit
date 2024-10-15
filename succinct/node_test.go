package succinct

import (
	"fmt"
	"testing"
)
import "github.com/stretchr/testify/assert"

// f(a) = x^2 + 5
func TestNodeBuilder_SimpleArithmetic(t *testing.T) {
	n := NewNodeBuilder()

	x := n.Var()
	xSq := n.Mul(x, x)
	five := n.Const(5)
	xSquaredPlusFive := n.Add(xSq, five)
	total := n.Add(xSquaredPlusFive, x)

	n.AssertEq(total, n.Val(11))

	n.FillNodes(map[int]float64{
		x.GetId(): 2,
	})
	assert.True(t, n.Verify(total))

	assert.Equal(t, float64(2), x.val)
	assert.Equal(t, float64(4), xSq.val)
	assert.Equal(t, float64(5), five.val)
	assert.Equal(t, float64(9), xSquaredPlusFive.val)
	assert.Equal(t, float64(11), total.val)
}

// f(a) = (a + 1) / 8
func TestNodeBuilder_DivisionHint(t *testing.T) {
	testTable := []struct {
		input  float64
		output float64
	}{
		{3, 0.5},
		{-3, -0.25},
		{63, 8},
	}

	for _, testCase := range testTable {
		t.Run(fmt.Sprintf("input: %v output %v", testCase.input, testCase.output), func(t *testing.T) {
			n := NewNodeBuilder()
			h := n.Hint()

			a := n.Var()
			one := n.Val(1)
			b := n.Add(a, one)
			c := h.Build(h.Div(h.Val(b), h.Const(8)))
			eight := n.Const(8)
			cTimesEight := n.Mul(c, eight)

			n.AssertEq(cTimesEight, b)
			n.AssertEq(c, n.Const(testCase.output))

			n.FillNodes(map[int]float64{
				a.GetId(): testCase.input,
			})
			assert.True(t, n.Verify(b))
			assert.Equal(t, testCase.input, a.GetVal())
			assert.Equal(t, testCase.output, c.GetVal())
		})
	}
}

// f(a) = sqrt(a + 7)
func TestNewNodeBuilder_SqrtHint(t *testing.T) {
	testTable := []struct {
		input  float64
		output float64
	}{
		{2, 3},
		{9, 4},
		{-7, 0},
	}

	for _, testCase := range testTable {
		t.Run(fmt.Sprintf("input: %v output %v", testCase.input, testCase.output), func(t *testing.T) {
			n := NewNodeBuilder()
			h := n.Hint()
			x := n.Var()
			seven := n.Const(7)
			xPlusSeven := n.Add(x, seven)

			y := h.Build(h.Sqrt(h.Val(xPlusSeven)))
			ySquared := n.Mul(y, y)

			n.AssertEq(xPlusSeven, ySquared)
			n.AssertEq(y, n.Const(testCase.output))

			n.FillNodes(map[int]float64{
				x.GetId(): testCase.input,
			})

			assert.True(t, n.Verify(xPlusSeven))
			assert.Equal(t, testCase.input, x.val)
			assert.Equal(t, testCase.output, y.val)
		})
	}
}

func TestNewNodeBuilder_Panics(t *testing.T) {
	n := NewNodeBuilder()
	h := n.h

	assert.Panicsf(t, func() {
		h.Div(h.Const(8), h.Const(0))
	}, "cannot divide by 0")
	assert.Panicsf(t, func() {
		h.Sqrt(h.Const(-23))
	}, "cannot take the square root of a negative number (at least in this circuit :))")
}

// f(a) = sqrt(x^2 + x + 2) - 4
func TestNewNodeBuilder_MultipleHints(t *testing.T) {
	testTable := []struct {
		input  float64
		output float64
	}{
		{1, -2},
		{-2, -2},
	}

	for _, testCase := range testTable {
		t.Run(fmt.Sprintf("input: %v output %v", testCase.input, testCase.output), func(t *testing.T) {
			n := NewNodeBuilder()
			h := n.Hint()

			// right side
			x := n.Var()
			xSquared := n.Mul(x, x)
			xSquaredPlusX := n.Add(xSquared, x)
			two := n.Val(2)
			xSquaredPlusXPlus2 := n.Add(xSquaredPlusX, two)
			firstTerm := h.Sqrt(h.Val(xSquaredPlusXPlus2))

			// left side
			y := h.Build(h.Sub(firstTerm, h.Const(4)))
			yPlusFour := n.Add(y, n.Const(4))

			n.AssertEq(yPlusFour, xSquaredPlusXPlus2)

			n.FillNodes(map[int]float64{
				x.GetId(): testCase.input,
			})
			n.Verify(xSquaredPlusXPlus2)
			assert.Equal(t, testCase.input, x.val)
			assert.Equal(t, testCase.output, y.val)
		})
	}
}

// f(a) = (x + y) / 4
func TestNewNodeBuilder_MultipleInputs(t *testing.T) {
	testTable := []struct {
		input1 float64
		input2 float64
		output float64
	}{
		{1, 5, 1.5},
		{-2, 72, 17.5},
		{-2, 2, 0},
	}

	for _, testCase := range testTable {
		t.Run(fmt.Sprintf("input1: %v, input2: %v, output %v", testCase.input1, testCase.input2, testCase.output),
			func(t *testing.T) {
				n := NewNodeBuilder()
				h := n.Hint()

				x := n.Var()
				y := n.Var()

				xPlusY := n.Add(x, y)
				xPlusYByFour := h.Build(h.Div(h.Val(xPlusY), h.Const(4)))
				n.AssertEq(n.Mul(xPlusYByFour, n.Const(4)), xPlusY)

				n.FillNodes(map[int]float64{
					x.GetId(): testCase.input1,
					y.GetId(): testCase.input2,
				})
				assert.True(t, n.Verify(xPlusY))

				assert.Equal(t, testCase.output, xPlusYByFour.val)
				assert.Equal(t, testCase.input1, x.val)
				assert.Equal(t, testCase.input2, y.val)
			})
	}
}
