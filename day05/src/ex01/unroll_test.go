package unroll

import (
	toytree "day05/ex00"
	"testing"
)

func TestUnroll(t *testing.T) {
	tree := toytree.TreeNode{
		HasToy: true,
		Left: &toytree.TreeNode{
			HasToy: true,
			Left: &toytree.TreeNode{
				HasToy: true,
				Left:   nil,
				Right:  nil,
			},
			Right: &toytree.TreeNode{
				HasToy: false,
				Left:   nil,
				Right:  nil,
			},
		},
		Right: &toytree.TreeNode{
			HasToy: false,
			Left: &toytree.TreeNode{
				HasToy: true,
				Left:   nil,
				Right:  nil,
			},
			Right: &toytree.TreeNode{
				HasToy: true,
				Left:   nil,
				Right:  nil,
			},
		},
	}

	got := unrollGarland(&tree)
	expect := []bool{true, true, false, true, true, false, true}
	for i, val := range got {
		if expect[i] != val {
			t.Error("Got not expected slice")
		}
	}
}
