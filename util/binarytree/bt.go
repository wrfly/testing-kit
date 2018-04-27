package binarytree

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// BuildTree by layers
func BuildTree(list []int) *TreeNode {
	root := &TreeNode{}
	queue := []*TreeNode{}
	currentNode := root
	for i := 0; i < len(list); i++ {
		// init the root
		if i == 0 {
			currentNode.Val = list[i]
			queue = append(queue, currentNode)
			i++
		}

		// next value of the list,
		// if value is valid, append it to the tree
		// as well as the queue
		if v := list[i]; v > 0 {
			currentNode.Left = &TreeNode{Val: v}
			queue = append(queue, currentNode.Left)
		}

		// i++ to get the next value
		// then check it
		if i++; i >= len(list) {
			break
		}

		// we got the next value of the list,
		// do the same for the right of the node
		if v := list[i]; v > 0 {
			currentNode.Right = &TreeNode{Val: v}
			queue = append(queue, currentNode.Right)
		}

		// we have finished this node,
		// pop it from the queue and re-point currentNode
		queue = queue[1:]
		currentNode = queue[0]
	}

	// we ranged over the list and got a tree
	return root
}
