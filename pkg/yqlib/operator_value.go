package yqlib

func valueOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	return context.SingleChildContext(expressionNode.Operation.CandidateNode), nil
}
