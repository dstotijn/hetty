package proj

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sync"
	"time"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/filter"
	"github.com/dstotijn/hetty/pkg/proxy/intercept"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/sender"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

// Service is used for managing projects.
type Service interface {
	CreateProject(ctx context.Context, name string) (Project, error)
	OpenProject(ctx context.Context, projectID ulid.ULID) (Project, error)
	CloseProject() error
	DeleteProject(ctx context.Context, projectID ulid.ULID) error
	ActiveProject(ctx context.Context) (Project, error)
	IsProjectActive(projectID ulid.ULID) bool
	Projects(ctx context.Context) ([]Project, error)
	Scope() *scope.Scope
	SetScopeRules(ctx context.Context, rules []scope.Rule) error
	SetRequestLogFindFilter(ctx context.Context, filter reqlog.FindRequestsFilter) error
	SetSenderRequestFindFilter(ctx context.Context, filter sender.FindRequestsFilter) error
	UpdateInterceptSettings(ctx context.Context, settings intercept.Settings) error
}

type service struct {
	repo            Repository
	interceptSvc    *intercept.Service
	reqLogSvc       reqlog.Service
	senderSvc       sender.Service
	scope           *scope.Scope
	activeProjectID ulid.ULID
	mu              sync.RWMutex
}

type Project struct {
	ID       ulid.ULID
	Name     string
	Settings Settings

	isActive bool
}

type Settings struct {
	// Request log settings
	ReqLogBypassOutOfScope bool
	ReqLogOnlyFindInScope  bool
	ReqLogSearchExpr       filter.Expression

	// Intercept settings
	InterceptRequests       bool
	InterceptResponses      bool
	InterceptRequestFilter  filter.Expression
	InterceptResponseFilter filter.Expression

	// Sender settings
	SenderOnlyFindInScope bool
	SenderSearchExpr      filter.Expression

	// Scope settings
	ScopeRules []scope.Rule
}

var (
	ErrProjectNotFound = errors.New("proj: project not found")
	ErrNoProject       = errors.New("proj: no open project")
	ErrNoSettings      = errors.New("proj: settings not found")
	ErrInvalidName     = errors.New("proj: invalid name, must be alphanumeric or whitespace chars")
)

var nameRegexp = regexp.MustCompile(`^[\w\d\s]+$`)

type Config struct {
	Repository       Repository
	InterceptService *intercept.Service
	ReqLogService    reqlog.Service
	SenderService    sender.Service
	Scope            *scope.Scope
}

// NewService returns a new Service.
func NewService(cfg Config) (Service, error) {
	return &service{
		repo:         cfg.Repository,
		interceptSvc: cfg.InterceptService,
		reqLogSvc:    cfg.ReqLogService,
		senderSvc:    cfg.SenderService,
		scope:        cfg.Scope,
	}, nil
}

func (svc *service) CreateProject(ctx context.Context, name string) (Project, error) {
	if !nameRegexp.MatchString(name) {
		return Project{}, ErrInvalidName
	}

	project := Project{
		ID:   ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
		Name: name,
	}

	err := svc.repo.UpsertProject(ctx, project)
	if err != nil {
		return Project{}, fmt.Errorf("proj: could not create project: %w", err)
	}

	return project, nil
}

// CloseProject closes the currently open project (if there is one).
func (svc *service) CloseProject() error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if svc.activeProjectID.Compare(ulid.ULID{}) == 0 {
		return nil
	}

	svc.activeProjectID = ulid.ULID{}
	svc.reqLogSvc.SetActiveProjectID(ulid.ULID{})
	svc.reqLogSvc.SetBypassOutOfScopeRequests(false)
	svc.reqLogSvc.SetFindReqsFilter(reqlog.FindRequestsFilter{})
	svc.interceptSvc.UpdateSettings(intercept.Settings{
		RequestsEnabled:  false,
		ResponsesEnabled: false,
		RequestFilter:    nil,
		ResponseFilter:   nil,
	})
	svc.senderSvc.SetActiveProjectID(ulid.ULID{})
	svc.senderSvc.SetFindReqsFilter(sender.FindRequestsFilter{})
	svc.scope.SetRules(nil)

	return nil
}

// DeleteProject removes a project from the repository.
func (svc *service) DeleteProject(ctx context.Context, projectID ulid.ULID) error {
	if svc.activeProjectID.Compare(projectID) == 0 {
		return fmt.Errorf("proj: project (%v) is active", projectID.String())
	}

	if err := svc.repo.DeleteProject(ctx, projectID); err != nil {
		return fmt.Errorf("proj: could not delete project: %w", err)
	}

	return nil
}

// OpenProject sets a project as the currently active project.
func (svc *service) OpenProject(ctx context.Context, projectID ulid.ULID) (Project, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	project, err := svc.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		return Project{}, fmt.Errorf("proj: failed to get project: %w", err)
	}

	svc.activeProjectID = project.ID

	// Request log settings.
	svc.reqLogSvc.SetFindReqsFilter(reqlog.FindRequestsFilter{
		ProjectID:   project.ID,
		OnlyInScope: project.Settings.ReqLogOnlyFindInScope,
		SearchExpr:  project.Settings.ReqLogSearchExpr,
	})
	svc.reqLogSvc.SetBypassOutOfScopeRequests(project.Settings.ReqLogBypassOutOfScope)
	svc.reqLogSvc.SetActiveProjectID(project.ID)

	// Intercept settings.
	svc.interceptSvc.UpdateSettings(intercept.Settings{
		RequestsEnabled:  project.Settings.InterceptRequests,
		ResponsesEnabled: project.Settings.InterceptResponses,
		RequestFilter:    project.Settings.InterceptRequestFilter,
		ResponseFilter:   project.Settings.InterceptResponseFilter,
	})

	// Sender settings.
	svc.senderSvc.SetActiveProjectID(project.ID)
	svc.senderSvc.SetFindReqsFilter(sender.FindRequestsFilter{
		ProjectID:   project.ID,
		OnlyInScope: project.Settings.SenderOnlyFindInScope,
		SearchExpr:  project.Settings.SenderSearchExpr,
	})

	// Scope settings.
	svc.scope.SetRules(project.Settings.ScopeRules)

	return project, nil
}

func (svc *service) ActiveProject(ctx context.Context) (Project, error) {
	activeProjectID := svc.activeProjectID
	if activeProjectID.Compare(ulid.ULID{}) == 0 {
		return Project{}, ErrNoProject
	}

	project, err := svc.repo.FindProjectByID(ctx, activeProjectID)
	if err != nil {
		return Project{}, fmt.Errorf("proj: failed to get active project: %w", err)
	}

	project.isActive = true

	return project, nil
}

func (svc *service) Projects(ctx context.Context) ([]Project, error) {
	projects, err := svc.repo.Projects(ctx)
	if err != nil {
		return nil, fmt.Errorf("proj: could not get projects: %w", err)
	}

	return projects, nil
}

func (svc *service) Scope() *scope.Scope {
	return svc.scope
}

func (svc *service) SetScopeRules(ctx context.Context, rules []scope.Rule) error {
	project, err := svc.ActiveProject(ctx)
	if err != nil {
		return err
	}

	project.Settings.ScopeRules = rules

	err = svc.repo.UpsertProject(ctx, project)
	if err != nil {
		return fmt.Errorf("proj: failed to update project: %w", err)
	}

	svc.scope.SetRules(rules)

	return nil
}

func (svc *service) SetRequestLogFindFilter(ctx context.Context, filter reqlog.FindRequestsFilter) error {
	project, err := svc.ActiveProject(ctx)
	if err != nil {
		return err
	}

	filter.ProjectID = project.ID

	project.Settings.ReqLogOnlyFindInScope = filter.OnlyInScope
	project.Settings.ReqLogSearchExpr = filter.SearchExpr

	err = svc.repo.UpsertProject(ctx, project)
	if err != nil {
		return fmt.Errorf("proj: failed to update project: %w", err)
	}

	svc.reqLogSvc.SetFindReqsFilter(filter)

	return nil
}

func (svc *service) SetSenderRequestFindFilter(ctx context.Context, filter sender.FindRequestsFilter) error {
	project, err := svc.ActiveProject(ctx)
	if err != nil {
		return err
	}

	filter.ProjectID = project.ID

	project.Settings.SenderOnlyFindInScope = filter.OnlyInScope
	project.Settings.SenderSearchExpr = filter.SearchExpr

	err = svc.repo.UpsertProject(ctx, project)
	if err != nil {
		return fmt.Errorf("proj: failed to update project: %w", err)
	}

	svc.senderSvc.SetFindReqsFilter(filter)

	return nil
}

func (svc *service) IsProjectActive(projectID ulid.ULID) bool {
	return projectID.Compare(svc.activeProjectID) == 0
}

func (svc *service) UpdateInterceptSettings(ctx context.Context, settings intercept.Settings) error {
	project, err := svc.ActiveProject(ctx)
	if err != nil {
		return err
	}

	project.Settings.InterceptRequests = settings.RequestsEnabled
	project.Settings.InterceptResponses = settings.ResponsesEnabled
	project.Settings.InterceptRequestFilter = settings.RequestFilter
	project.Settings.InterceptResponseFilter = settings.ResponseFilter

	err = svc.repo.UpsertProject(ctx, project)
	if err != nil {
		return fmt.Errorf("proj: failed to update project: %w", err)
	}

	svc.interceptSvc.UpdateSettings(settings)

	return nil
}
