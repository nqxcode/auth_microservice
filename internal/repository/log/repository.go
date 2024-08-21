package log

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	"github.com/nqxcode/platform_common/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "log"

	idColumn        = "log_id"
	messageColumn   = "message"
	payloadColumn   = "payload"
	ipColumn        = "ip"
	createdAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository new log repository
func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, model *model.Log) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(messageColumn, payloadColumn, ipColumn).
		Values(model.Message, model.Payload, model.IP).
		Suffix("RETURNING " + idColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     tableName + "_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
