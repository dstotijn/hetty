package proj

import (
	"context"

	"github.com/oklog/ulid"
)

type Repository interface {
	FindProjectByID(ctx context.Context, id ulid.ULID) (Project, error)
	UpsertProject(ctx context.Context, project Project) error
	DeleteProject(ctx context.Context, id ulid.ULID) error
	Projects(ctx context.Context) ([]Project, error)
	Close() error
}
