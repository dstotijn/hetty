package intercept

import "github.com/dstotijn/hetty/pkg/search"

type Settings struct {
	Enabled       bool
	RequestFilter search.Expression
}
