package search

import "strings"

type Expression interface {
	String() string
}

type PrefixExpression struct {
	Operator TokenType
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	b := strings.Builder{}
	b.WriteString("(")
	b.WriteString(pe.Operator.String())
	b.WriteString(" ")
	b.WriteString(pe.Right.String())
	b.WriteString(")")

	return b.String()
}

type InfixExpression struct {
	Operator TokenType
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) String() string {
	b := strings.Builder{}
	b.WriteString("(")
	b.WriteString(ie.Left.String())
	b.WriteString(" ")
	b.WriteString(ie.Operator.String())
	b.WriteString(" ")
	b.WriteString(ie.Right.String())
	b.WriteString(")")

	return b.String()
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) String() string {
	return sl.Value
}
