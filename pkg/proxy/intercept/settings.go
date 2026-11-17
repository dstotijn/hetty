package intercept

import "github.com/pocikode/ratatoskr/pkg/filter"

type Settings struct {
	RequestsEnabled  bool
	ResponsesEnabled bool
	RequestFilter    filter.Expression
	ResponseFilter   filter.Expression
}
