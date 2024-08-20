package user

import (
	"context"
	"fmt"

	"github.com/nqxcode/auth_microservice/internal/client/db"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/repository"
	"github.com/nqxcode/auth_microservice/internal/repository/user/converter"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/model"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "user"

	idColumn        = "user_id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, model *model.User) (int64, error) {
	builder := sq.Insert(escape(tableName)).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(model.Info.Name, model.Info.Email, model.Info.Role, model.Password).
		Suffix("RETURNING " + idColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     tableName + "_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.UpdateUserInfo) error {
	if info == nil {
		return nil
	}

	builder := sq.Update(escape(tableName)).
		PlaceholderFormat(sq.Dollar)

	if info.Name != nil {
		builder = builder.Set(nameColumn, *info.Name)
	}
	if info.Role != nil {
		builder = builder.Set(roleColumn, *info.Role)
	}

	builder.
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     tableName + "_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(escape(tableName)).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     tableName + "_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Find(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, passwordColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(escape(tableName)).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     tableName + "_repository.Find",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func escape(value string) string {
	return fmt.Sprintf("\"%s\"", value)
}
