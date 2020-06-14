package main

import (
	"github.com/somen440/gonv/converter"
	"github.com/somen440/gonv/structure"
)

// Operator operate ask
type Operator struct {
	before *structure.DatabaseStructure
	after  *structure.DatabaseStructure
}

// NewOperator return Operator
func NewOperator(before, after *structure.DatabaseStructure) *Operator {
	return &Operator{
		before: before,
		after:  after,
	}
}

// Ask questions and generate answers
func (o *Operator) Ask() *converter.ModifiedAnswer {
	// todo: operate

	return nil
}
