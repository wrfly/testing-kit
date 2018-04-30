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

func preorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	r := []int{root.Val}
	r = append(r, preorderTraversal(root.Left)...)
	return append(r, preorderTraversal(root.Right)...)
}

func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	r := inorderTraversal(root.Left)
	r = append(r, root.Val)
	return append(r, inorderTraversal(root.Right)...)
}

func postorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	r := postorderTraversal(root.Left)
	r = append(r, postorderTraversal(root.Right)...)
	return append(r, root.Val)
}

func inorderTraversal_2(root *TreeNode) []int {
	if root == nil {
		return nil
	}

	stack := []*TreeNode{root}
	r := []int{}
	leftVisited := false

	for len(stack) != 0 {
		// left first
		if root.Left != nil && !leftVisited {
			root = root.Left
			// append this node to stack
			stack = append(stack, root)
			continue
		}

		// left is nil, append this node
		r = append(r, root.Val)
		stack = stack[:len(stack)-1]

		// if right is not nil, point root to it
		if root.Right != nil {
			root = root.Right
			// append this node to stack
			stack = append(stack, root)
			// reset the `leftVisited` flag
			leftVisited = false
			continue
		}

		// when we reached the last one, stack is empty
		// break here
		if len(stack) == 0 {
			break
		}
		// pop one node from stack and set the `leftVisited` flag
		root = stack[len(stack)-1]
		leftVisited = true
	}
	return r
}

func inorderTraversal_3(root *TreeNode) []int {
	if root == nil {
		return nil
	}

	stack := []*TreeNode{}
	r := []int{}

	for {
		stack = append(stack, root)
		if root.Left != nil {
			root = root.Left
			continue
		}

	pop:
		root = stack[len(stack)-1]
		r = append(r, root.Val)
		stack = stack[:len(stack)-1]
		if root.Right != nil {
			root = root.Right
			continue
		}
		if len(stack) == 0 {
			break
		}
		goto pop
	}
	return r
}
