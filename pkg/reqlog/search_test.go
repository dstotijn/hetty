package reqlog_test

import (
	"testing"

	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/search"
)

func TestRequestLogMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		query         string
		requestLog    reqlog.RequestLog
		expectedMatch bool
		expectedError error
	}{
		{
			name:  "infix expression, equal operator, match",
			query: "req.body = foo",
			requestLog: reqlog.RequestLog{
				Body: []byte("foo"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, not equal operator, match",
			query: "req.body != bar",
			requestLog: reqlog.RequestLog{
				Body: []byte("foo"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, greater than operator, match",
			query: "req.body > a",
			requestLog: reqlog.RequestLog{
				Body: []byte("b"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, less than operator, match",
			query: "req.body < b",
			requestLog: reqlog.RequestLog{
				Body: []byte("a"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, greater than or equal operator, match greater than",
			query: "req.body >= a",
			requestLog: reqlog.RequestLog{
				Body: []byte("b"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, greater than or equal operator, match equal",
			query: "req.body >= a",
			requestLog: reqlog.RequestLog{
				Body: []byte("a"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, less than or equal operator, match less than",
			query: "req.body <= b",
			requestLog: reqlog.RequestLog{
				Body: []byte("a"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, less than or equal operator, match equal",
			query: "req.body <= b",
			requestLog: reqlog.RequestLog{
				Body: []byte("b"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, regular expression operator, match",
			query: `req.body =~ "^foo(.*)$"`,
			requestLog: reqlog.RequestLog{
				Body: []byte("foobar"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, negate regular expression operator, match",
			query: `req.body !~ "^foo(.*)$"`,
			requestLog: reqlog.RequestLog{
				Body: []byte("xoobar"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, and operator, match",
			query: "req.body = bar AND res.body = yolo",
			requestLog: reqlog.RequestLog{
				Body: []byte("bar"),
				Response: &reqlog.ResponseLog{
					Body: []byte("yolo"),
				},
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "infix expression, or operator, match",
			query: "req.body = bar OR res.body = yolo",
			requestLog: reqlog.RequestLog{
				Body: []byte("foo"),
				Response: &reqlog.ResponseLog{
					Body: []byte("yolo"),
				},
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "prefix expression, not operator, match",
			query: "NOT (req.body = bar)",
			requestLog: reqlog.RequestLog{
				Body: []byte("foo"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "string literal expression, match in request log",
			query: "foo",
			requestLog: reqlog.RequestLog{
				Body: []byte("foo"),
			},
			expectedMatch: true,
			expectedError: nil,
		},
		{
			name:  "string literal expression, no match",
			query: "foo",
			requestLog: reqlog.RequestLog{
				Body: []byte("bar"),
			},
			expectedMatch: false,
			expectedError: nil,
		},
		{
			name:  "string literal expression, match in response log",
			query: "foo",
			requestLog: reqlog.RequestLog{
				Response: &reqlog.ResponseLog{
					Body: []byte("foo"),
				},
			},
			expectedMatch: true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			searchExpr, err := search.ParseQuery(tt.query)
			assertError(t, nil, err)

			got, err := tt.requestLog.Matches(searchExpr)
			assertError(t, tt.expectedError, err)

			if tt.expectedMatch != got {
				t.Errorf("expected match result: %v, got: %v", tt.expectedMatch, got)
			}
		})
	}
}

func assertError(t *testing.T, exp, got error) {
	t.Helper()

	switch {
	case exp == nil && got != nil:
		t.Fatalf("expected: nil, got: %v", got)
	case exp != nil && got == nil:
		t.Fatalf("expected: %v, got: nil", exp.Error())
	case exp != nil && got != nil && exp.Error() != got.Error():
		t.Fatalf("expected: %v, got: %v", exp.Error(), got.Error())
	}
}
