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
	x_squared_plus_five := n.Add(x_sq, five)
	total := n.Add(x_squared_plus_five, x)

	y := n.Var()
	n.AssertEq(total, y)

	n.FillNodes(map[int]int{
		x.GetId(): 2,
		y.GetId(): 11,
	})
	fmt.Println(n.Verify(total))

}
