package proj

import (
	"context"
)

type Repository interface {
	UpsertSettings(ctx context.Context, module string, settings interface{}) error
	FindSettingsByModule(ctx context.Context, module string, settings interface{}) error
	OpenProject(name string) error
	DeleteProject(name string) error
	Projects() ([]Project, error)
	Close() error
}
