package yqlib

import (
	"fmt"
)

type DataTreeNavigator interface {
	// given the context and a expressionNode,
	// this will process the against the given expressionNode and return
	// a new context of matching candidates
	GetMatchingNodes(context Context, expressionNode *ExpressionNode) (Context, error)
}

type dataTreeNavigator struct {
}

func NewDataTreeNavigator() DataTreeNavigator {
	return &dataTreeNavigator{}
}

func (d *dataTreeNavigator) GetMatchingNodes(context Context, expressionNode *ExpressionNode) (Context, error) {
	if expressionNode == nil {
		return context, nil
	}
	handler := expressionNode.Operation.OperationType.Handler
	if handler != nil {
		return handler(d, context, expressionNode)
	}
	return Context{}, fmt.Errorf("Unknown operator %v", expressionNode.Operation.OperationType)

}
