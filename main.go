package main

import (
	"awesomeProject/succinct"
	"fmt"
)

func main() {
	n := succinct.NewNodeBuilder()

	x := n.Var()

	x_sq := n.Mul(x, x)
	five := n.Const(5)
	first_two := n.Add(x_sq, x)

	y := n.Var()
	n.AssertEq(y, first_two)

	total := n.Add(first_two, five)
	n.FillNodes(map[int]int{
		x.GetId(): 2,
		y.GetId(): 11,
	})
	fmt.Println(n.Verify(total))
}
