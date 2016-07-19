// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package vm

import "fmt"

// a Visitor's VisitBefore method is invoked for each node encountered by Walk.
// If the result Visitor v is not nil, Walk visits each of the children of that
// node with v.  VisitAfter is called on n at the end.
type Visitor interface {
	VisitBefore(n node) (v Visitor)
	VisitAfter(n node)
}

// convenience function
func walknodelist(v Visitor, list []node) {
	for _, x := range list {
		Walk(v, x)
	}
}

func Walk(v Visitor, node node) {
	// Returning nil from VisitBefore signals to Walk that the Visitor has
	// handled the children of this node.
	if v := v.VisitBefore(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *stmtlistNode:
		walknodelist(v, n.children)
	case *exprlistNode:
		walknodelist(v, n.children)

	case *condNode:
		if n.cond != nil {
			Walk(v, n.cond)
		}
		Walk(v, n.truthNode)
		if n.elseNode != nil {
			Walk(v, n.elseNode)
		}

	case *builtinNode:
		if n.args != nil {
			Walk(v, n.args)
		}

	case *binaryExprNode:
		Walk(v, n.lhs)
		Walk(v, n.rhs)

	case *unaryExprNode:
		Walk(v, n.lhs)

	case *indexedExprNode:
		Walk(v, n.index)
		Walk(v, n.lhs)

	case *defNode:
		walknodelist(v, n.children)

	case *decoNode:
		walknodelist(v, n.children)

	case *regexNode, *idNode, *caprefNode, *declNode, *stringConstNode, *intConstNode, *floatConstNode, *nextNode, *otherwiseNode:
		// These nodes are terminals, thus have no children to walk.

	default:
		panic(fmt.Sprintf("Walk: unexpected node type %T: %v", n, n))
	}

	v.VisitAfter(node)
}
