package scope

import "context"

type Repository interface {
	UpsertSettings(ctx context.Context, module string, settings interface{}) error
	FindSettingsByModule(ctx context.Context, module string, settings interface{}) error
}
