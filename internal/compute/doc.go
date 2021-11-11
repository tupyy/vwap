package compute

/*

This package computes the volume average of points.

stack.go implements a LILO stack

vwap.go

It computes the volume average. The computation is not O(n) but O(1). It keeps the total volume of all points in the stack as variable
as well the sum of products value*volume for all points in the stack. If the stack is full, when a new point is added, it substract the popped point
from total volume and the sum of product and adds the new one keeping in this way the two variable consistent with the content of the stack.
*/
