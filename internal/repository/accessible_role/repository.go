package log

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/nqxcode/platform_common/client/db"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	"github.com/nqxcode/auth_microservice/internal/repository/accessible_role/converter"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/accessible_role/model"
)

const (
	tableName = "accessible_role"

	idColumn              = "accessible_role_id"
	roleColumn            = "role"
	endpointAddressColumn = "endpoint_address"
	createdAtColumn       = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository new accessible role repository
func NewRepository(db db.Client) repository.AccessibleRoleRepository {
	return &repo{db: db}
}

// GetList get list of accessible roles
func (r *repo) GetList(ctx context.Context) ([]model.AccessibleRole, error) {
	builder := sq.Select(idColumn, roleColumn, endpointAddressColumn, createdAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(escape(tableName))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     tableName + "_repository.GetList",
		QueryRaw: query,
	}

	var roles []modelRepo.AccessibleRole
	err = r.db.DB().ScanAllContext(ctx, &roles, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToManyAccessibleRoleFromRepo(roles), nil
}

func escape(value string) string {
	return fmt.Sprintf("\"%s\"", value)
}
