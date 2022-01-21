package api

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/oklog/ulid"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/search"
)

type Resolver struct {
	ProjectService    proj.Service
	RequestLogService *reqlog.Service
}

type (
	queryResolver    struct{ *Resolver }
	mutationResolver struct{ *Resolver }
)

func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *queryResolver) HTTPRequestLogs(ctx context.Context) ([]HTTPRequestLog, error) {
	reqs, err := r.RequestLogService.FindRequests(ctx)
	if errors.Is(err, proj.ErrNoProject) {
		return nil, noActiveProjectErr(ctx)
	} else if err != nil {
		return nil, fmt.Errorf("could not query repository for requests: %w", err)
	}

	logs := make([]HTTPRequestLog, len(reqs))

	for i, req := range reqs {
		req, err := parseRequestLog(req)
		if err != nil {
			return nil, err
		}

		logs[i] = req
	}

	return logs, nil
}

func (r *queryResolver) HTTPRequestLog(ctx context.Context, id ULID) (*HTTPRequestLog, error) {
	log, err := r.RequestLogService.FindRequestLogByID(ctx, ulid.ULID(id))
	if errors.Is(err, reqlog.ErrRequestNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not get request by ID: %w", err)
	}

	req, err := parseRequestLog(log)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func parseRequestLog(reqLog reqlog.RequestLog) (HTTPRequestLog, error) {
	method := HTTPMethod(reqLog.Method)
	if method != "" && !method.IsValid() {
		return HTTPRequestLog{}, fmt.Errorf("request has invalid method: %v", method)
	}

	log := HTTPRequestLog{
		ID:        ULID(reqLog.ID),
		Proto:     reqLog.Proto,
		Method:    method,
		Timestamp: ulid.Time(reqLog.ID.Time()),
	}

	if reqLog.URL != nil {
		log.URL = reqLog.URL.String()
	}

	if len(reqLog.Body) > 0 {
		bodyStr := string(reqLog.Body)
		log.Body = &bodyStr
	}

	if reqLog.Header != nil {
		log.Headers = make([]HTTPHeader, 0)

		for key, values := range reqLog.Header {
			for _, value := range values {
				log.Headers = append(log.Headers, HTTPHeader{
					Key:   key,
					Value: value,
				})
			}
		}
	}

	if reqLog.Response != nil {
		log.Response = &HTTPResponseLog{
			Proto:      reqLog.Response.Proto,
			StatusCode: reqLog.Response.StatusCode,
		}
		statusReasonSubs := strings.SplitN(reqLog.Response.Status, " ", 2)

		if len(statusReasonSubs) == 2 {
			log.Response.StatusReason = statusReasonSubs[1]
		}

		if len(reqLog.Response.Body) > 0 {
			bodyStr := string(reqLog.Response.Body)
			log.Response.Body = &bodyStr
		}

		if reqLog.Response.Header != nil {
			log.Response.Headers = make([]HTTPHeader, 0)

			for key, values := range reqLog.Response.Header {
				for _, value := range values {
					log.Response.Headers = append(log.Response.Headers, HTTPHeader{
						Key:   key,
						Value: value,
					})
				}
			}
		}
	}

	return log, nil
}

func (r *mutationResolver) CreateProject(ctx context.Context, name string) (*Project, error) {
	p, err := r.ProjectService.CreateProject(ctx, name)
	if errors.Is(err, proj.ErrInvalidName) {
		return nil, gqlerror.Errorf("Project name must only contain alphanumeric or space chars.")
	} else if err != nil {
		return nil, fmt.Errorf("could not open project: %w", err)
	}

	return &Project{
		ID:       ULID(p.ID),
		Name:     p.Name,
		IsActive: r.ProjectService.IsProjectActive(p.ID),
	}, nil
}

func (r *mutationResolver) OpenProject(ctx context.Context, id ULID) (*Project, error) {
	p, err := r.ProjectService.OpenProject(ctx, ulid.ULID(id))
	if errors.Is(err, proj.ErrInvalidName) {
		return nil, gqlerror.Errorf("Project name must only contain alphanumeric or space chars.")
	} else if err != nil {
		return nil, fmt.Errorf("could not open project: %w", err)
	}

	return &Project{
		ID:       ULID(p.ID),
		Name:     p.Name,
		IsActive: r.ProjectService.IsProjectActive(p.ID),
	}, nil
}

func (r *queryResolver) ActiveProject(ctx context.Context) (*Project, error) {
	p, err := r.ProjectService.ActiveProject(ctx)
	if errors.Is(err, proj.ErrNoProject) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("could not open project: %w", err)
	}

	return &Project{
		ID:       ULID(p.ID),
		Name:     p.Name,
		IsActive: r.ProjectService.IsProjectActive(p.ID),
	}, nil
}

func (r *queryResolver) Projects(ctx context.Context) ([]Project, error) {
	p, err := r.ProjectService.Projects(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get projects: %w", err)
	}

	projects := make([]Project, len(p))
	for i, proj := range p {
		projects[i] = Project{
			ID:       ULID(proj.ID),
			Name:     proj.Name,
			IsActive: r.ProjectService.IsProjectActive(proj.ID),
		}
	}

	return projects, nil
}

func (r *queryResolver) Scope(ctx context.Context) ([]ScopeRule, error) {
	rules := r.ProjectService.Scope().Rules()
	return scopeToScopeRules(rules), nil
}

func regexpToStringPtr(r *regexp.Regexp) *string {
	if r == nil {
		return nil
	}

	s := r.String()

	return &s
}

func (r *mutationResolver) CloseProject(ctx context.Context) (*CloseProjectResult, error) {
	if err := r.ProjectService.CloseProject(); err != nil {
		return nil, fmt.Errorf("could not close project: %w", err)
	}

	return &CloseProjectResult{true}, nil
}

func (r *mutationResolver) DeleteProject(ctx context.Context, id ULID) (*DeleteProjectResult, error) {
	if err := r.ProjectService.DeleteProject(ctx, ulid.ULID(id)); err != nil {
		return nil, fmt.Errorf("could not delete project: %w", err)
	}

	return &DeleteProjectResult{
		Success: true,
	}, nil
}

func (r *mutationResolver) ClearHTTPRequestLog(ctx context.Context) (*ClearHTTPRequestLogResult, error) {
	project, err := r.ProjectService.ActiveProject(ctx)
	if errors.Is(err, proj.ErrNoProject) {
		return nil, noActiveProjectErr(ctx)
	} else if err != nil {
		return nil, fmt.Errorf("could not get active project: %w", err)
	}

	if err := r.RequestLogService.ClearRequests(ctx, project.ID); err != nil {
		return nil, fmt.Errorf("could not clear request log: %w", err)
	}

	return &ClearHTTPRequestLogResult{true}, nil
}

func (r *mutationResolver) SetScope(ctx context.Context, input []ScopeRuleInput) ([]ScopeRule, error) {
	rules := make([]scope.Rule, len(input))

	for i, rule := range input {
		u, err := stringPtrToRegexp(rule.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL in scope rule: %w", err)
		}

		var headerKey, headerValue *regexp.Regexp

		if rule.Header != nil {
			headerKey, err = stringPtrToRegexp(rule.Header.Key)
			if err != nil {
				return nil, fmt.Errorf("invalid header key in scope rule: %w", err)
			}

			headerValue, err = stringPtrToRegexp(rule.Header.Key)
			if err != nil {
				return nil, fmt.Errorf("invalid header value in scope rule: %w", err)
			}
		}

		body, err := stringPtrToRegexp(rule.Body)
		if err != nil {
			return nil, fmt.Errorf("invalid body in scope rule: %w", err)
		}

		rules[i] = scope.Rule{
			URL: u,
			Header: scope.Header{
				Key:   headerKey,
				Value: headerValue,
			},
			Body: body,
		}
	}

	err := r.ProjectService.SetScopeRules(ctx, rules)
	if err != nil {
		return nil, fmt.Errorf("could not set scope rules: %w", err)
	}

	return scopeToScopeRules(rules), nil
}

func (r *queryResolver) HTTPRequestLogFilter(ctx context.Context) (*HTTPRequestLogFilter, error) {
	return findReqFilterToHTTPReqLogFilter(r.RequestLogService.FindReqsFilter), nil
}

func (r *mutationResolver) SetHTTPRequestLogFilter(
	ctx context.Context,
	input *HTTPRequestLogFilterInput,
) (*HTTPRequestLogFilter, error) {
	filter, err := findRequestsFilterFromInput(input)
	if err != nil {
		return nil, fmt.Errorf("could not parse request log filter: %w", err)
	}

	err = r.ProjectService.SetRequestLogFindFilter(ctx, filter)
	if errors.Is(err, proj.ErrNoProject) {
		return nil, noActiveProjectErr(ctx)
	} else if err != nil {
		return nil, fmt.Errorf("could not set request log filter: %w", err)
	}

	return findReqFilterToHTTPReqLogFilter(filter), nil
}

func stringPtrToRegexp(s *string) (*regexp.Regexp, error) {
	if s == nil {
		return nil, nil
	}

	return regexp.Compile(*s)
}

func scopeToScopeRules(rules []scope.Rule) []ScopeRule {
	scopeRules := make([]ScopeRule, len(rules))
	for i, rule := range rules {
		scopeRules[i].URL = regexpToStringPtr(rule.URL)
		if rule.Header.Key != nil || rule.Header.Value != nil {
			scopeRules[i].Header = &ScopeHeader{
				Key:   regexpToStringPtr(rule.Header.Key),
				Value: regexpToStringPtr(rule.Header.Value),
			}
		}

		scopeRules[i].Body = regexpToStringPtr(rule.Body)
	}

	return scopeRules
}

func findRequestsFilterFromInput(input *HTTPRequestLogFilterInput) (filter reqlog.FindRequestsFilter, err error) {
	if input == nil {
		return
	}

	if input.OnlyInScope != nil {
		filter.OnlyInScope = *input.OnlyInScope
	}

	if input.SearchExpression != nil && *input.SearchExpression != "" {
		expr, err := search.ParseQuery(*input.SearchExpression)
		if err != nil {
			return reqlog.FindRequestsFilter{}, fmt.Errorf("could not parse search query: %w", err)
		}

		filter.SearchExpr = expr
	}

	return
}

func findReqFilterToHTTPReqLogFilter(findReqFilter reqlog.FindRequestsFilter) *HTTPRequestLogFilter {
	empty := reqlog.FindRequestsFilter{}
	if findReqFilter == empty {
		return nil
	}

	httpReqLogFilter := &HTTPRequestLogFilter{
		OnlyInScope: findReqFilter.OnlyInScope,
	}

	if findReqFilter.SearchExpr != nil {
		searchExpr := findReqFilter.SearchExpr.String()
		httpReqLogFilter.SearchExpression = &searchExpr
	}

	return httpReqLogFilter
}

func noActiveProjectErr(ctx context.Context) error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "No active project.",
		Extensions: map[string]interface{}{
			"code": "no_active_project",
		},
	}
}
