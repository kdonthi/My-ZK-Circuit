package main

import (
	"awesomeProject/succinct"
	"fmt"
)

func main() {
	//n := succinct.NewNodeBuilder()
	//
	//x := n.Var()
	//x_sq := n.Mul(x, x)
	//five := n.Const(5)
	//x_squared_plus_five := n.Add(x_sq, five)
	//total := n.Add(x_squared_plus_five, x)
	//
	//y := n.Var()
	//n.AssertEq(total, y)
	//
	//n.FillNodes(map[int]float64{
	//	x.GetId(): 2,
	//	y.GetId(): 11,
	//})
	//fmt.Println(n.Verify(total))

	//Example 2: f(a) = (a+1) / 8
	//
	//function f(a):
	//    b = a + 1
	//    c = b / 8
	//    return c

	n2 := succinct.NewNodeBuilder()
	h2 := n2.Hint()

	a := n2.Var()
	one := n2.Num(1)
	b := n2.Add(a, one)

	c := h2.Build(h2.Div(h2.Val(b), h2.Const(8))) // get a node through the n2!
	eight := n2.Const(8)
	c_times_eight := n2.Mul(c, eight)

	n2.AssertEq(c_times_eight, b)

	n2.FillNodes(map[int]float64{
		a.GetId(): 1,
		c.GetId(): 0.25,
	})
	fmt.Println(n2.Verify(b))
	fmt.Println(c.GetVal(), a.GetVal(), b.GetVal())
	// TODO add one more to the end?
}
