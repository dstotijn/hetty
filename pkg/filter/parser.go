package filter

import (
	"fmt"
	"regexp"
)

type precedence int

const (
	_ precedence = iota
	precLowest
	precAnd
	precOr
	precNot
	precEq
	precLessGreater
	precPrefix
	precGroup
)

type (
	prefixParser func(*Parser) (Expression, error)
	infixParser  func(*Parser, Expression) (Expression, error)
)

var (
	prefixParsers = map[TokenType]prefixParser{}
	infixParsers  = map[TokenType]infixParser{}
)

var tokenPrecedences = map[TokenType]precedence{
	TokParenOpen: precGroup,
	TokOpNot:     precNot,
	TokOpAnd:     precAnd,
	TokOpOr:      precOr,
	TokOpEq:      precEq,
	TokOpNotEq:   precEq,
	TokOpGt:      precLessGreater,
	TokOpLt:      precLessGreater,
	TokOpGtEq:    precLessGreater,
	TokOpLtEq:    precLessGreater,
	TokOpRe:      precEq,
	TokOpNotRe:   precEq,
}

func init() {
	// Populate maps in `init`, because package global variables would cause an
	// initialization cycle.
	infixOperators := []TokenType{
		TokOpAnd,
		TokOpOr,
		TokOpEq,
		TokOpNotEq,
		TokOpGt,
		TokOpLt,
		TokOpGtEq,
		TokOpLtEq,
		TokOpRe,
		TokOpNotRe,
	}
	for _, op := range infixOperators {
		infixParsers[op] = parseInfixExpression
	}

	prefixParsers[TokOpNot] = parsePrefixExpression
	prefixParsers[TokString] = parseStringLiteral
	prefixParsers[TokParenOpen] = parseGroupedExpression
}

// maxParseDepth bounds the recursion depth of the parser. Filter expressions
// are short, human-written search queries, so any legitimate query nests far
// below this limit. Without a bound, deeply nested input (e.g. many "(" or
// repeated "NOT" prefixes) recurses until the goroutine stack is exhausted,
// which is a *fatal*, unrecoverable runtime error in Go (it cannot be caught by
// recover, so it crashes the whole process). ParseQuery is reachable from the
// unauthenticated GraphQL API via search expressions, so this would be a remote
// denial of service.
const maxParseDepth = 256

var ErrMaxDepthExceeded = fmt.Errorf("filter: expression nesting exceeds maximum depth of %d", maxParseDepth)

type Parser struct {
	l     *Lexer
	cur   Token
	peek  Token
	depth int
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	return p
}

func ParseQuery(input string) (expr Expression, err error) {
	p := &Parser{l: NewLexer(input)}
	// Ensure the lexer goroutine is always released, including on early returns
	// caused by parse errors that leave unread tokens on the channel.
	defer p.l.Stop()

	p.nextToken()
	p.nextToken()

	if p.curTokenIs(TokEOF) {
		return nil, fmt.Errorf("filter: unexpected EOF")
	}

	for !p.curTokenIs(TokEOF) {
		right, err := p.parseExpression(precLowest)

		switch {
		case err != nil:
			return nil, fmt.Errorf("filter: could not parse expression: %w", err)
		case expr == nil:
			expr = right
		default:
			expr = InfixExpression{
				Operator: TokOpAnd,
				Left:     expr,
				Right:    right,
			}
		}

		p.nextToken()
	}

	return
}

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.Next()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.cur.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peek.Type == t
}

func (p *Parser) curPrecedence() precedence {
	if p, ok := tokenPrecedences[p.cur.Type]; ok {
		return p
	}

	return precLowest
}

func (p *Parser) peekPrecedence() precedence {
	if p, ok := tokenPrecedences[p.peek.Type]; ok {
		return p
	}

	return precLowest
}

func (p *Parser) parseExpression(prec precedence) (Expression, error) {
	// Guard against stack exhaustion from deeply nested input. parseExpression
	// is the single recursion point reached by grouped, prefix and infix
	// sub-expression parsing, so bounding it here bounds total recursion depth.
	p.depth++
	defer func() { p.depth-- }()

	if p.depth > maxParseDepth {
		return nil, ErrMaxDepthExceeded
	}

	prefixParser, ok := prefixParsers[p.cur.Type]
	if !ok {
		return nil, fmt.Errorf("no prefix parse function for %v found", p.cur.Type)
	}

	expr, err := prefixParser(p)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression prefix: %w", err)
	}

	for !p.peekTokenIs(eof) && prec < p.peekPrecedence() {
		infixParser, ok := infixParsers[p.peek.Type]
		if !ok {
			break
		}

		p.nextToken()

		expr, err = infixParser(p, expr)
		if err != nil {
			return nil, fmt.Errorf("could not parse infix expression: %w", err)
		}
	}

	return expr, nil
}

func parsePrefixExpression(p *Parser) (Expression, error) {
	expr := PrefixExpression{
		Operator: p.cur.Type,
	}

	p.nextToken()

	right, err := p.parseExpression(precPrefix)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression for right operand: %w", err)
	}

	expr.Right = right

	return expr, nil
}

func parseInfixExpression(p *Parser, left Expression) (Expression, error) {
	expr := InfixExpression{
		Operator: p.cur.Type,
		Left:     left,
	}

	prec := p.curPrecedence()
	p.nextToken()

	right, err := p.parseExpression(prec)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression for right operand: %w", err)
	}

	if expr.Operator == TokOpRe || expr.Operator == TokOpNotRe {
		if rightStr, ok := right.(StringLiteral); ok {
			re, err := regexp.Compile(rightStr.Value)
			if err != nil {
				return nil, fmt.Errorf("could not compile regular expression %q: %w", rightStr.Value, err)
			}

			right = RegexpLiteral{re}
		}
	}

	expr.Right = right

	return expr, nil
}

func parseStringLiteral(p *Parser) (Expression, error) {
	return StringLiteral{Value: p.cur.Literal}, nil
}

func parseGroupedExpression(p *Parser) (Expression, error) {
	p.nextToken()

	expr, err := p.parseExpression(precLowest)
	if err != nil {
		return nil, fmt.Errorf("could not parse grouped expression: %w", err)
	}

	for p.nextToken(); !p.curTokenIs(TokParenClose); p.nextToken() {
		if p.curTokenIs(TokEOF) {
			return nil, fmt.Errorf("unexpected EOF: unmatched parentheses")
		}

		right, err := p.parseExpression(precLowest)
		if err != nil {
			return nil, fmt.Errorf("could not parse expression: %w", err)
		}

		expr = InfixExpression{
			Operator: TokOpAnd,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}
