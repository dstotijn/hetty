package reqlog

import (
	"unicode"

	"github.com/db47h/lex"
	"github.com/db47h/lex/state"
)

const (
	tokEOF = iota
	tokString
	tokOpNot
	tokOpAnd
	tokOpOr
	tokOpEq
	tokOpNeq
	tokOpGt
	tokOpGteq
	tokOpLt
	tokOpLteq
	tokOpHas
	tokOpRe
	tokOpNre
	tokParenOpen
	tokParenClose
)

type lexItem struct {
	token lex.Token
	value string
}

func lexQuery(s *lex.State) lex.StateFn {
	str := lexString()
	quotedString := state.QuotedString(tokString)

	return func(s *lex.State) lex.StateFn {
		r := s.Next()
		pos := s.Pos()
		switch r {
		case lex.EOF:
			s.Emit(pos, tokEOF, nil)
			return nil
		case '"':
			return quotedString
		case '=':
			if next := s.Next(); next == '~' {
				s.Emit(pos, tokOpRe, nil)
			} else {
				s.Backup()
				s.Emit(pos, tokOpEq, nil)
			}
			return nil
		case '!':
			switch next := s.Next(); next {
			case '=':
				s.Emit(pos, tokOpNeq, nil)
				return nil
			case '~':
				s.Emit(pos, tokOpNre, nil)
				return nil
			default:
				s.Backup()
			}
		case '>':
			if next := s.Next(); next == '=' {
				s.Emit(pos, tokOpGteq, nil)
			} else {
				s.Backup()
				s.Emit(pos, tokOpGt, nil)
			}
			return nil
		case '<':
			if next := s.Next(); next == '=' {
				s.Emit(pos, tokOpLteq, nil)
			} else {
				s.Backup()
				s.Emit(pos, tokOpLt, nil)
			}
			return nil
		case ':':
			s.Emit(pos, tokOpHas, nil)
			return nil
		case '(':
			s.Emit(pos, tokParenOpen, nil)
			return nil
		case ')':
			s.Emit(pos, tokParenClose, nil)
			return nil
		}

		switch {
		case unicode.IsSpace(r):
			// Absorb spaces.
			for r = s.Next(); unicode.IsSpace(r); r = s.Next() {
			}
			s.Backup()
			return nil
		default:
			return str
		}
	}
}

func lexString() lex.StateFn {
	// Preallocate a buffer to store the value. It will end-up being at
	// least as large as the largest value scanned.
	b := make([]rune, 0, 64)

	isStringChar := func(r rune) bool {
		switch r {
		case '=', '!', '<', '>', ':', '(', ')':
			return false
		}
		return !(unicode.IsSpace(r) || r == lex.EOF)
	}

	return func(l *lex.State) lex.StateFn {
		pos := l.Pos()
		// Reset buffer and add first char.
		b = append(b[:0], l.Current())
		// Read identifier.
		for r := l.Next(); isStringChar(r); r = l.Next() {
			b = append(b, r)
		}
		// The character returned by the last call to `l.Next` is not part of
		// the value. Undo it.
		l.Backup()

		switch {
		case string(b) == "NOT":
			l.Emit(pos, tokOpNot, nil)
		case string(b) == "AND":
			l.Emit(pos, tokOpAnd, nil)
		case string(b) == "OR":
			l.Emit(pos, tokOpOr, nil)
		default:
			l.Emit(pos, tokString, string(b))
		}

		return nil
	}
}
