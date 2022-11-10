package yqlib

import "container/list"

func unionOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	lhs, err := d.GetMatchingNodes(context, expressionNode.LHS)
	if err != nil {
		return Context{}, err
	}
	rhs, err := d.GetMatchingNodes(context, expressionNode.RHS)

	if err != nil {
		return Context{}, err
	}

	results := lhs.ChildContext(list.New())
	for el := lhs.MatchingNodes.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		results.MatchingNodes.PushBack(node)
	}

	// this can happen when both expressions modify the context
	// instead of creating their own.
	/// (.foo = "bar"), (.thing = "cat")
	if rhs.MatchingNodes != lhs.MatchingNodes {

		for el := rhs.MatchingNodes.Front(); el != nil; el = el.Next() {
			node := el.Value.(*CandidateNode)

			results.MatchingNodes.PushBack(node)
		}
	}
	return results, nil
}
