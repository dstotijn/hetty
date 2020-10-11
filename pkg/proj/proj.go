package proj

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dstotijn/hetty/pkg/db/sqlite"
	"github.com/dstotijn/hetty/pkg/scope"
)

// Service is used for managing projects.
type Service struct {
	dbPath string
	db     *sqlite.Client
	name   string

	Scope *scope.Scope
}

type Project struct {
	Name     string
	IsActive bool
}

var (
	ErrNoProject   = errors.New("proj: no open project")
	ErrInvalidName = errors.New("proj: invalid name, must be alphanumeric or whitespace chars")
)

var nameRegexp = regexp.MustCompile(`^[\w\d\s]+$`)

// NewService returns a new Service.
func NewService(dbPath string) (*Service, error) {
	// Create directory for DBs if it doesn't exist yet.
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dbPath, 0755); err != nil {
			return nil, fmt.Errorf("proj: could not create project directory: %v", err)
		}
	}

	return &Service{
		dbPath: dbPath,
		db:     &sqlite.Client{},
		Scope:  scope.New(nil),
	}, nil
}

// Close closes the currently open project database (if there is one).
func (svc *Service) Close() error {
	if err := svc.db.Close(); err != nil {
		return fmt.Errorf("proj: could not close project: %v", err)
	}
	svc.name = ""
	return nil
}

// Delete removes a project database file from disk (if there is one).
func (svc *Service) Delete(name string) error {
	if name == "" {
		return errors.New("proj: name cannot be empty")
	}
	if svc.name == name {
		return fmt.Errorf("proj: project (%v) is active", name)
	}

	if err := os.Remove(filepath.Join(svc.dbPath, name+".db")); err != nil {
		return fmt.Errorf("proj: could not remove database file: %v", err)
	}

	return nil
}

// Database returns the currently open database. If no database is open, it will
// return `nil`.
func (svc *Service) Database() *sqlite.Client {
	return svc.db
}

// Open opens a database identified with `name`. If a database with this
// identifier doesn't exist yet, it will be automatically created.
func (svc *Service) Open(name string) (Project, error) {
	if !nameRegexp.MatchString(name) {
		return Project{}, ErrInvalidName
	}
	if err := svc.db.Close(); err != nil {
		return Project{}, fmt.Errorf("proj: could not close previously open database: %v", err)
	}

	dbPath := filepath.Join(svc.dbPath, name+".db")

	err := svc.db.Open(dbPath)
	if err != nil {
		return Project{}, fmt.Errorf("proj: could not open database: %v", err)
	}

	svc.name = name

	return Project{
		Name:     name,
		IsActive: true,
	}, nil
}

func (svc *Service) ActiveProject() (Project, error) {
	if !svc.db.IsOpen() {
		return Project{}, ErrNoProject
	}

	return Project{
		Name: svc.name,
	}, nil
}

func (svc *Service) Projects() ([]Project, error) {
	files, err := ioutil.ReadDir(svc.dbPath)
	if err != nil {
		return nil, fmt.Errorf("proj: could not read projects directory: %v", err)
	}

	projects := make([]Project, len(files))
	for i, file := range files {
		projName := strings.TrimSuffix(file.Name(), ".db")
		projects[i] = Project{
			Name:     projName,
			IsActive: svc.name == projName,
		}
	}

	return projects, nil
}
