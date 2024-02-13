package dblib

import (
	"context"
	"database/sql"
	"errors"

	"htmx/model"
	"htmx/pkg/dbtx"
	"htmx/pkg/util"
)

func Get[T model.GetTableName](ctx context.Context, query string, args ...any) (util.Null[T], error) {
	tx := dbtx.GetReadonly(ctx)
	var result T
	err := tx.GetContext(ctx, &result, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return util.Null[T]{}, nil
		}
		return util.Null[T]{}, err
	}
	return util.Null[T]{
		Valid: true,
		Data:  result,
	}, nil
}

func Insert[T ~int64](ctx context.Context, query string, data any) (T, error) {
	tx := dbtx.GetTx(ctx)
	result, err := tx.NamedExecContext(ctx, query, data)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return T(id), nil
}
