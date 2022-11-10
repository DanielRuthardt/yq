package yqlib

import "gopkg.in/yaml.v3"

func equalsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, isEquals(false), true)
}

func isEquals(flip bool) func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
	return func(d *dataTreeNavigator, context Context, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
		value := false
		if lhs == nil && rhs == nil {
			owner := &CandidateNode{}
			return createBooleanCandidate(owner, !flip), nil
		} else if lhs == nil {
			rhsNode := unwrapDoc(rhs.Node)
			value := rhsNode.Tag == "!!null"
			if flip {
				value = !value
			}
			return createBooleanCandidate(rhs, value), nil
		} else if rhs == nil {
			lhsNode := unwrapDoc(lhs.Node)
			value := lhsNode.Tag == "!!null"
			if flip {
				value = !value
			}
			return createBooleanCandidate(lhs, value), nil
		}

		lhsNode := unwrapDoc(lhs.Node)
		rhsNode := unwrapDoc(rhs.Node)

		if lhsNode.Tag == "!!null" {
			value = (rhsNode.Tag == "!!null")
		} else if lhsNode.Kind == yaml.ScalarNode && rhsNode.Kind == yaml.ScalarNode {
			value = matchKey(lhsNode.Value, rhsNode.Value)
		}
		if flip {
			value = !value
		}
		return createBooleanCandidate(lhs, value), nil
	}
}

func notEqualsOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return crossFunction(d, context.ReadOnlyClone(), expressionNode, isEquals(true), true)
}
