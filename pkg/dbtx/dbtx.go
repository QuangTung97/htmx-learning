package dbtx

import (
	"context"
	"database/sql"

	"github.com/QuangTung97/svloc"
	"github.com/jmoiron/sqlx"

	"htmx/config"
)

type Readonly interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Transaction interface {
	Readonly

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

var _ Transaction = &sqlx.DB{}
var _ Transaction = &sqlx.Tx{}

// Provider for creating Readonly and Transaction
type Provider interface {
	Transact(ctx context.Context, fn func(ctx context.Context) error) error
	Readonly(ctx context.Context) context.Context

	// Autocommit only for testing and special cases
	Autocommit(ctx context.Context) context.Context
}

type providerImpl struct {
	db *sqlx.DB
}

var ProviderLoc = svloc.Register[Provider](func(unv *svloc.Universe) Provider {
	db := config.DBLoc.Get(unv)
	return &providerImpl{db: db}
})

func (p *providerImpl) Transact(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	_, ok := getTxFromContext(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	ctx = context.WithValue(ctx, ctxTxKey, ctxTxValue{
		tx: tx,
	})

	err = fn(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *providerImpl) Readonly(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxReadonlyKey, ctxReadonlyValue{
		db: p.db,
	})
}

func (p *providerImpl) Autocommit(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxTxKey, ctxTxValue{
		tx: p.db,
	})
}
