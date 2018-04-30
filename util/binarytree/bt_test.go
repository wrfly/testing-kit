package binarytree

import (
	"testing"
)

/*
tree := []int{1, -1, 2, 3, -1, 7, -1, -1, 4}

    1
-1     2
     3  -1
   7  -1
-1   4

*/

func TestBinaryTree(t *testing.T) {
	tree := []int{1, -1, 2, 3, -1, 7, -1, -1, 4}

	t.Run("Build tree", func(t *testing.T) {
		root := BuildTree(tree)
		t.Log(root.Right.Left.Left.Val)
	})

	t.Run("Traversal", func(t *testing.T) {
		root := BuildTree(tree)
		t.Log(preorderTraversal(root))
		t.Log(inorderTraversal(root))
		t.Log(inorderTraversal_2(root))
		t.Log(inorderTraversal_3(root))
		t.Log(postorderTraversal(root))
	})
}
