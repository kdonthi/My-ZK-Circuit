# ZK Circuit

This library allows users to create and verify a ZK circuit involving addition and multiplication, but also utilizing subtraction, square roots, and division using hints.

This document will try to explain how to use the library to be able to verify equations (like you have probably been wanting to your whole life!)

Let's start with something really simple:
`f(a) = x + 5`
```go
nb := succinct.NodeBuilder() // creates the builder instance

x := nb.Var() // creates a variable you have to set later
five := nb.Val(5) // creates a constant
x_plus_five := nb.Add(x, five)

nb.AssertEq(x_plus_five, n2.Const(6))
nb.FillNodes(map[int]float64{
    x.GetId(): 1, // sets the variable you created earlier
})

var verified bool
verified = nb.Verify(x_plus_five) // to verify, put in the head of your node chain
```

As you can see, we just had to create the variables we wanted, and chain them together with others to create the equation. 

Now let's see how to use these hints I was talking about earlier:

`f(a) = sqrt(x + 4)`
```go
nb := succinct.NodeBuilder()
h := nb.Hint()

x := nb.Var()
four := nb.Val(4)
x_plus_four := nb.Add(x, four)

hint := h.Sqrt(h.Val(x_plus_four)) // use h.Val() to convert a Node to a HintNode
y := h.Build(hint) // use h.Build() to convert a HintNode to a Node

y_squared := nb.Mul(y, y)

nb.AssertEq(y_squared, x_plus_four)
nb.FillNodes(map[int]float64{
    x.GetId(): 1,
})

var verified bool
verified = nb.Verify(x_plus_four) // the head of the node chain with no hints/equality is x_plus_four (hints and equality are computed separately)
```

There are more tests under `succinct/node_test.go` if you are interested in some more complex examples. What to put inside 
`Verify` might be slightly confusing, but just try to remember to put the head of the node chain that is not based on any hints because
the hints and the underlying logic are calculated seperately.




