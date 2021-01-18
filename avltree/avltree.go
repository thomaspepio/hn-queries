package avltree

type NodeType string

type rebalancingStrategy string

const (
	Root       NodeType = "root"
	LeftChild  NodeType = "leftChild"
	RightChild NodeType = "rightChild"

	rightRight    rebalancingStrategy = "rightright" // insert at the right side, balance factor >= 2
	leftleft      rebalancingStrategy = "leftleft"   // insert at the left side, balance factor <= -2
	rightleft     rebalancingStrategy = "rightleft"  // insert at the right side, balance factor <= -2
	leftright     rebalancingStrategy = "leftright"  // insert at the left side, balance factor >= 2
	noRebalancing rebalancingStrategy = "none"
)

// An AVLTree whose keys are fixed to integers
type AVLTree struct {
	Key      int
	Values   map[int]int
	Left     *AVLTree
	Right    *AVLTree
	Parent   *AVLTree
	NodeType NodeType
}

// New returns leafless tree, with height set to 0
func New(key int, values map[int]int) *AVLTree {
	return &AVLTree{key, values, nil, nil, nil, Root}
}

func newLeftTree(key int, values map[int]int, parent *AVLTree) *AVLTree {
	return &AVLTree{key, values, nil, nil, parent, LeftChild}
}

func newRightTree(key int, values map[int]int, parent *AVLTree) *AVLTree {
	return &AVLTree{key, values, nil, nil, parent, RightChild}
}

// Get : lookup a key in the tree
func (tree *AVLTree) Get(key int) map[int]int {
	if tree != nil {
		if key == tree.Key {
			return tree.Values
		} else if key < tree.Key {
			return tree.Left.Get(key)
		} else if key > tree.Key {
			return tree.Right.Get(key)
		}
	}

	return nil
}

// Update : when the key is present, replaces it's associated value
func (tree *AVLTree) Update(key int, values map[int]int) {
	if tree != nil {
		if key == tree.Key {
			tree.Values = values
		} else if key < tree.Key {
			tree.Left.Update(key, values)
		} else if key > tree.Key {
			tree.Right.Update(key, values)
		}
	}
}

// Insert : self balancing insertion
func (tree *AVLTree) Insert(key int, values map[int]int) {
	if key < tree.Key {
		if tree.Left == nil {
			tree.Left = newLeftTree(key, values, tree)
		} else {
			tree.Left.Insert(key, values)
		}
	} else if key > tree.Key {
		if tree.Right == nil {
			tree.Right = newRightTree(key, values, tree)
		} else {
			tree.Right.Insert(key, values)
		}
	}

	tree.balance()
}

// Height : computes the height of a tree
// A leaf has height == 0
// A Tree with one leaf | right | both children has height == 1
func (tree *AVLTree) Height() int {
	if tree == nil {
		return -1
	}

	return 1 + max(tree.Left.Height(), tree.Right.Height())
}

// Balance = Height(right subtree) - Height(left subtree)
// Balance == 0 means the tree is balanced
// Balance < 0 indicates a "left heavy" tree
// Balance > 0 indicates a "right heavy" tree
// AVL trees maintain the following invariant : Balance(root) ∈ {-1, 0, 1}
func (tree *AVLTree) Balance() int {
	if tree == nil {
		return 0
	}

	return tree.Right.Height() - tree.Left.Height()
}

// LeftRotate : rotation of an assumed right heavy tree (balance factor > 0)
//	From : A			To : 	B
//	 		\				  /	  \
//	  		 B				 A	   C
//	   		   \
//		        C
func (tree *AVLTree) LeftRotate() {
	if !(tree == nil || tree.Right == nil) {
		newRoot := *tree.Right
		formerRoot := *tree

		newRoot.Parent = formerRoot.Parent
		newRoot.NodeType = formerRoot.NodeType

		formerRoot.NodeType = LeftChild
		formerRoot.Right = newRoot.Left

		newRoot.Left = &formerRoot

		*tree = newRoot
		groomLeft(tree)
		groomRight(tree)
	}
}

// RightRotate : rotation of an assumed left heavy tree (balance factor < 0)
//	From :   C			To : 	B
//	 		/				  /	  \
//	  	   B				 A	   C
//	   	  /
//	     A
func (tree *AVLTree) RightRotate() {
	if !(tree == nil || tree.Left == nil) {
		newRoot := *tree.Left
		formerRoot := *tree

		newRoot.Parent = formerRoot.Parent
		newRoot.NodeType = formerRoot.NodeType

		formerRoot.NodeType = RightChild
		formerRoot.Left = newRoot.Right

		newRoot.Right = &formerRoot

		*tree = newRoot
		groomLeft(tree)
		groomRight(tree)
	}
}

// LeftRightRotate : left child of the root with balancing factor > 0
//	From :   C			To : 	B
//	 		/				  /	  \
//	  	   A				 A	   C
//	   	    \
//	         B
func (tree *AVLTree) LeftRightRotate() {
	if !(tree == nil || tree.Left == nil) {
		tree.Left.LeftRotate()
		tree.RightRotate()
	}
}

// RightLeftRotate : right child of the root with balancing factor < 0
//	From :   A			To : 	B
//	 		  \ 			  /	  \
//	  	       C    		 A	   C
//	   	      /
//	         B
func (tree *AVLTree) RightLeftRotate() {
	if !(tree == nil || tree.Right == nil) {
		tree.Right.RightRotate()
		tree.LeftRotate()
	}
}

func (tree *AVLTree) balance() {
	parentBalance := tree.Parent.Balance()
	if parentBalance < -1 || parentBalance > 1 {
		rebalanceStrategy := getRebalanceStrategy(tree)
		switch rebalanceStrategy {
		case rightRight:
			tree.Parent.LeftRotate()
		case leftleft:
			tree.Parent.RightRotate()
		case rightleft:
			tree.Parent.RightLeftRotate()
		case leftright:
			tree.Parent.LeftRightRotate()
		}
	}
}

func getRebalanceStrategy(tree *AVLTree) rebalancingStrategy {
	balance := tree.Balance()
	nodeType := tree.NodeType

	if nodeType == RightChild {
		if balance >= 0 {
			return rightRight
		} else if balance < 0 {
			return rightleft
		}
	} else if nodeType == LeftChild {
		if balance > 0 {
			return leftright
		} else if balance <= 0 {
			return leftleft
		}
	}

	return noRebalancing
}

func groomLeft(tree *AVLTree) {
	if tree.Left != nil {
		tree.Left.Parent = tree
		tree.Left.NodeType = LeftChild

		leftLeft := tree.Left.Left
		if leftLeft != nil {
			leftLeft.Parent = tree.Left
		}

		leftRight := tree.Left.Right
		if leftRight != nil {
			leftRight.Parent = tree.Left
		}
	}
}

func groomRight(tree *AVLTree) {
	if tree.Right != nil {
		tree.Right.Parent = tree
		tree.Right.NodeType = RightChild

		rightLeft := tree.Right.Left
		if rightLeft != nil {
			rightLeft.Parent = tree.Right
		}

		rightRight := tree.Right.Right
		if rightRight != nil {
			rightRight.Parent = tree.Right
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
