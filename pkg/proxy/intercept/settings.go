package intercept

import "github.com/dstotijn/hetty/pkg/search"

type Settings struct {
	RequestsEnabled  bool
	ResponsesEnabled bool
	RequestFilter    search.Expression
	ResponseFilter   search.Expression
}
