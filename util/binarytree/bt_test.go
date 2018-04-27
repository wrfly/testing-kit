package binarytree

import (
	"testing"
)

func TestBinaryTree(t *testing.T) {
	t.Run("build tree", func(t *testing.T) {
		list := []int{1, -1, 2, 3, -1, 7}
		root := BuildTree(list)
		t.Log(root.Right.Left.Left.Val)
	})

	t.Run("Traversal", func(t *testing.T) {
		list := []int{1, -1, 2, 3, -1, 7}
		root := BuildTree(list)
		t.Log(preorderTraversal(root))
		t.Log(inorderTraversal(root))
		t.Log(postorderTraversal(root))
	})
}
