package search

import (
	"encoding/gob"
	"regexp"
	"strconv"
	"strings"
)

type Expression interface {
	String() string
}

type PrefixExpression struct {
	Operator TokenType
	Right    Expression
}

func (pe PrefixExpression) String() string {
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

func (ie InfixExpression) String() string {
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

func (sl StringLiteral) String() string {
	return strconv.Quote(sl.Value)
}

type RegexpLiteral struct {
	*regexp.Regexp
}

func (rl RegexpLiteral) String() string {
	return strconv.Quote(rl.Regexp.String())
}

func (rl RegexpLiteral) MarshalBinary() ([]byte, error) {
	return []byte(rl.Regexp.String()), nil
}

func (rl *RegexpLiteral) UnmarshalBinary(data []byte) error {
	re, err := regexp.Compile(string(data))
	if err != nil {
		return err
	}

	*rl = RegexpLiteral{re}

	return nil
}

func init() {
	gob.Register(PrefixExpression{})
	gob.Register(InfixExpression{})
	gob.Register(StringLiteral{})
	gob.Register(RegexpLiteral{})
}
