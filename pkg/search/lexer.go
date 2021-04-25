package search

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

const eof = 0

// Token types.
const (
	// Flow.
	TokInvalid TokenType = iota
	TokEOF
	TokParenOpen
	TokParenClose

	// Literals.
	TokString

	// Boolean operators.
	TokOpNot
	TokOpAnd
	TokOpOr

	// Comparison operators.
	TokOpEq
	TokOpNotEq
	TokOpGt
	TokOpLt
	TokOpGtEq
	TokOpLtEq
	TokOpRe
	TokOpNotRe
)

var (
	keywords = map[string]TokenType{
		"NOT": TokOpNot,
		"AND": TokOpAnd,
		"OR":  TokOpOr,
	}
	reservedRunes    = []rune{'=', '!', '<', '>', '(', ')'}
	tokenTypeStrings = map[TokenType]string{
		TokInvalid:    "INVALID",
		TokEOF:        "EOF",
		TokParenOpen:  "(",
		TokParenClose: ")",
		TokString:     "STRING",
		TokOpNot:      "NOT",
		TokOpAnd:      "AND",
		TokOpOr:       "OR",
		TokOpEq:       "=",
		TokOpNotEq:    "!=",
		TokOpGt:       ">",
		TokOpLt:       "<",
		TokOpGtEq:     ">=",
		TokOpLtEq:     "<=",
		TokOpRe:       "=~",
		TokOpNotRe:    "!~",
	}
)

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input  string
	pos    int
	start  int
	width  int
	tokens chan Token
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}

	go l.run(begin)

	return l
}

func (l *Lexer) Next() Token {
	return <-l.tokens
}

func (tt TokenType) String() string {
	if typeString, ok := tokenTypeStrings[tt]; ok {
		return typeString
	}

	return "<unknown>"
}

func (l *Lexer) run(init stateFn) {
	for nextState := init; nextState != nil; {
		nextState = nextState(l)
	}
	close(l.tokens)
}

func (l *Lexer) read() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return
}

func (l *Lexer) emit(tokenType TokenType) {
	l.tokens <- Token{
		Type:    tokenType,
		Literal: l.input[l.start:l.pos],
	}

	l.start = l.pos
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) skip() {
	l.pos += l.width
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{
		Type:    TokInvalid,
		Literal: fmt.Sprintf(format, args...),
	}

	return nil
}

func begin(l *Lexer) stateFn {
	r := l.read()
	switch r {
	case '=':
		if next := l.read(); next == '~' {
			l.emit(TokOpRe)
		} else {
			l.backup()
			l.emit(TokOpEq)
		}

		return begin
	case '!':
		switch next := l.read(); next {
		case '=':
			l.emit(TokOpNotEq)
		case '~':
			l.emit(TokOpNotRe)
		default:
			return l.errorf("invalid rune %v", r)
		}

		return begin
	case '<':
		if next := l.read(); next == '=' {
			l.emit(TokOpLtEq)
		} else {
			l.backup()
			l.emit(TokOpLt)
		}

		return begin
	case '>':
		if next := l.read(); next == '=' {
			l.emit(TokOpGtEq)
		} else {
			l.backup()
			l.emit(TokOpGt)
		}

		return begin
	case '(':
		l.emit(TokParenOpen)
		return begin
	case ')':
		l.emit(TokParenClose)
		return begin
	case '"':
		return l.delimString(r)
	case eof:
		l.emit(TokEOF)
		return nil
	}

	if unicode.IsSpace(r) {
		l.ignore()
		return begin
	}

	return unquotedString
}

func (l *Lexer) delimString(delim rune) stateFn {
	// Ignore the start delimiter rune.
	l.ignore()

	for r := l.read(); r != delim; r = l.read() {
		if r == eof {
			return l.errorf("unexpected EOF, unclosed delimiter")
		}
	}
	// Don't include the end delimiter in emitted token.
	l.backup()
	l.emit(TokString)
	// Skip end delimiter.
	l.skip()

	return begin
}

func unquotedString(l *Lexer) stateFn {
	for r := l.read(); ; r = l.read() {
		switch {
		case r == eof:
			l.backup()
			l.emitUnquotedString()

			return begin
		case unicode.IsSpace(r):
			l.backup()
			l.emitUnquotedString()
			l.skip()

			return begin
		case isReserved(r):
			l.backup()
			l.emitUnquotedString()

			return begin
		}
	}
}

func (l *Lexer) emitUnquotedString() {
	str := l.input[l.start:l.pos]
	if tokType, ok := keywords[str]; ok {
		l.emit(tokType)
		return
	}

	l.emit(TokString)
}

func isReserved(r rune) bool {
	for _, v := range reservedRunes {
		if r == v {
			return true
		}
	}

	return false
}
