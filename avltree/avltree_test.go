package avltree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New_ShouldHave_NoBranches_NoParent_HeightOfZero_BalanceOfZero_NodeTypeRoot(t *testing.T) {
	tree := New(0, mapOf(0))
	assert.Nil(t, tree.Left, "Left branch should be nil")
	assert.Nil(t, tree.Right, "Right branch should be nil")
	assert.Nil(t, tree.Parent, "Parent should be nil")
	assert.Equal(t, 0, tree.Height(), "Height should be 0")
	assert.Equal(t, 0, tree.Balance(), "Balance should be 0")
	assert.Equal(t, Root, tree.NodeType, "Root should be Root node type")
}

func Test_Get_ElementExists(t *testing.T) {
	tree := New(0, mapOf(0))
	tree.Insert(1, mapOf(1))
	tree.Insert(2, mapOf(2))

	assert.Equal(t, mapOf(0), tree.Get(0), "Should have found \"0\"")
	assert.Equal(t, mapOf(1), tree.Get(1), "Should have found \"1\"")
	assert.Equal(t, mapOf(2), tree.Get(2), "Should have found \"2\"")
}

func Test_Get_ElementAbsent(t *testing.T) {
	tree := New(0, mapOf(0))
	assert.Nil(t, tree.Get(1), "Nil should be returned")
}

func Test_GetBetween_ShouldExclude_OutOfBoundsValues(t *testing.T) {
	tree := New(0, mapOf(0))
	tree.Insert(1, mapOf(1))
	tree.Insert(-1, mapOf(-1))
	tree.Insert(2, mapOf(2))
	tree.Insert(3, mapOf(3))
	tree.Insert(4, mapOf(4))
	tree.Insert(50, mapOf(50))

	actual := tree.Between(0, 5)

	assert.Equal(t, mapOf(0), actual.Get(0), "Resulting tree should contain key 0")
	assert.Equal(t, mapOf(1), actual.Get(1), "Resulting tree should contain key 1")
	assert.Equal(t, mapOf(2), actual.Get(2), "Resulting tree should contain key 2")
	assert.Equal(t, mapOf(3), actual.Get(3), "Resulting tree should contain key 3")
}

func Test_GetBetween_ShouldInclude_InBoundValues(t *testing.T) {
	tree := New(0, mapOf(0))
	tree.Insert(-1, mapOf(-1))
	tree.Insert(1, mapOf(1))
	tree.Insert(2, mapOf(2))
	tree.Insert(3, mapOf(3))
	tree.Insert(4, mapOf(4))

	actual := tree.Between(-1, 4)

	assert.Equal(t, mapOf(-1), actual.Get(-1), "Resulting tree should contain key -1")
	assert.Equal(t, mapOf(0), actual.Get(0), "Resulting tree should contain key 0")
	assert.Equal(t, mapOf(1), actual.Get(1), "Resulting tree should contain key 1")
	assert.Equal(t, mapOf(2), actual.Get(2), "Resulting tree should contain key 2")
	assert.Equal(t, mapOf(3), actual.Get(3), "Resulting tree should contain key 3")
	assert.Equal(t, mapOf(4), actual.Get(4), "Resulting tree should contain key 3")
}

func Test_Update_KeyExists(t *testing.T) {
	tree := New(0, mapOf(0))
	tree.Insert(1, mapOf(1))
	tree.Insert(2, mapOf(2))

	tree.Update(0, mapOf(10))
	tree.Update(1, mapOf(20))
	tree.Update(2, mapOf(30))

	assert.Equal(t, mapOf(10), tree.Get(0), "Should have found \"10\"")
	assert.Equal(t, mapOf(20), tree.Get(1), "Should have found \"20\"")
	assert.Equal(t, mapOf(30), tree.Get(2), "Should have found \"30\"")
}

func Test_RightInserts_ShouldIncrementHeight(t *testing.T) {
	tree := New(50, mapOf(0))
	tree.Insert(60, mapOf(1))
	assert.Equal(t, 1, tree.Height(), "Height should have been incremented")
}

func Test_LefThenRightInserts_ShouldKeepHeightSteady(t *testing.T) {
	tree := New(50, mapOf(0))
	tree.Insert(25, mapOf(1))
	tree.Insert(75, mapOf(2))
	assert.Equal(t, 1, tree.Height(), "Height should've been kept steady by left-right insert")
}

func Test_NodeWithOnlyLeftChild_ShouldHaveBalanceOfMinusOne(t *testing.T) {
	tree := New(50, mapOf(0))
	tree.Insert(25, mapOf(0))
	assert.Equal(t, -1, tree.Balance(), "A node with only a left child should have a balance of -1")
}

func Test_NodeWithOnlyRightChild_ShouldHaveBalanceOfMinusOne(t *testing.T) {
	tree := New(50, mapOf(0))
	tree.Insert(75, mapOf(0))
	assert.Equal(t, 1, tree.Balance(), "A node with only a right child should have a balance of 1")
}

func Test_Left_Rotate(t *testing.T) {
	expected := New(60, mapOf(0))
	expected.Insert(50, mapOf(0))
	expected.Insert(70, mapOf(0))

	actual := New(50, mapOf(0))
	actual.Insert(60, mapOf(0))
	actual.Insert(70, mapOf(0))

	leftRightAssertions(t, expected, actual)
	assert.Equal(t, expected, actual, "Left rotation is broken")
}

func Test_Right_Rotate(t *testing.T) {
	expected := New(60, mapOf(0))
	expected.Insert(50, mapOf(0))
	expected.Insert(70, mapOf(0))

	actual := New(70, mapOf(0))
	actual.Insert(60, mapOf(0))
	actual.Insert(50, mapOf(0))

	leftRightAssertions(t, expected, actual)
	assert.Equal(t, expected, actual, "Right rotation is broken")
}

func Test_LeftRight_Rotate(t *testing.T) {
	expected := New(60, mapOf(0))
	expected.Insert(50, mapOf(0))
	expected.Insert(70, mapOf(0))

	actual := New(70, mapOf(0))
	actual.Insert(50, mapOf(0))
	actual.Insert(60, mapOf(0))

	leftRightAssertions(t, expected, actual)
	assert.Equal(t, expected, actual, "Left-Right rotation is broken")
}

func Test_RightLeft_Rotate(t *testing.T) {
	expected := New(60, mapOf(0))
	expected.Insert(50, mapOf(0))
	expected.Insert(70, mapOf(0))

	actual := New(70, mapOf(0))
	actual.Insert(50, mapOf(0))
	actual.Insert(60, mapOf(0))

	leftRightAssertions(t, expected, actual)
	assert.Equal(t, expected, actual, "Right-Left rotation is broken !")
}

func Test_AVL_Invariant(t *testing.T) {
	tree := getTree(1000, false)
	assert.True(t, tree.Balance() == -1 || tree.Balance() == 0 || tree.Balance() == 1, "AVL invariant broken : balance should always be -1, 0 or 1")
}

func Test_ParentChild_SanityCheck(t *testing.T) {
	tree := getTree(1000, false)
	assert.True(t, parentChildSanityCheck(tree), "At least one left or right child does not references its parent")
}

func Test_NodeType_SanityCheck(t *testing.T) {
	tree := getTree(1000, false)
	assert.True(t, parentChildSanityCheck(tree), "At least one left or right node is mislabelled")
}

func getTree(n int, trace bool) *AVLTree {
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(n)

	tree := New(-1, mapOf(0))
	for _, r := range p[:n] {
		if trace {
			fmt.Println("Inserting ", r)
		}
		tree.Insert(r, mapOf(0))
	}

	return tree
}

func parentChildSanityCheck(tree *AVLTree) bool {
	if tree == nil {
		return true
	}

	leftChildReferencesParent := true
	if tree.Left != nil {
		leftChildReferencesParent = tree == tree.Left.Parent
	}

	rightChildReferencesParent := true
	if tree.Right != nil {
		rightChildReferencesParent = tree == tree.Right.Parent
	}

	return (leftChildReferencesParent && rightChildReferencesParent) && parentChildSanityCheck(tree.Left) && parentChildSanityCheck(tree.Right)
}

func nodeTypeSanityChekc(tree *AVLTree) bool {
	if tree == nil {
		return true
	}

	leftHasCorrectType := true
	if tree.Left != nil {
		leftHasCorrectType = tree.Left.NodeType == LeftChild
	}

	rightHasCorrectType := true
	if tree.Right != nil {
		rightHasCorrectType = tree.Right.NodeType == RightChild
	}

	return (leftHasCorrectType && rightHasCorrectType) && nodeTypeSanityChekc(tree.Left) && nodeTypeSanityChekc(tree.Right)
}

func leftRightAssertions(t *testing.T, expected *AVLTree, actual *AVLTree) {
	assert.Equal(t, expected.Key, actual.Key, "Root keys should be equal")
	assert.Equal(t, expected.Values, actual.Values, "Root values should be equal")
	assert.Equal(t, expected.Left.Key, actual.Left.Key, "Left keys shouls be equal")
	assert.Equal(t, expected.Left.Values, actual.Left.Values, "Left values should be equal")
	assert.Equal(t, expected.Right.Key, actual.Right.Key, "Right keys shouls be equal")
	assert.Equal(t, expected.Right.Values, actual.Right.Values, "Right values should be equal")
	assert.Equal(t, actual, actual.Left.Parent, "Left parent link broken")
	assert.Equal(t, actual, actual.Right.Parent, "Right parent link broken")
}

func mapOf(key int) map[int]int {
	return map[int]int{key: 0}
}
