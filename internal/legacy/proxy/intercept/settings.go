package intercept

import "github.com/pocikode/ratatoskr/internal/legacy/filter"

type Settings struct {
	RequestsEnabled  bool
	ResponsesEnabled bool
	RequestFilter    filter.Expression
	ResponseFilter   filter.Expression
}
