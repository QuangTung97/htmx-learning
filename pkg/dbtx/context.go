package dbtx

import (
	"context"
)

type ctxTxKeyType struct {
}

type ctxReadonlyKeyType struct {
}

var ctxTxKey = ctxTxKeyType{}
var ctxReadonlyKey = ctxReadonlyKeyType{}

type ctxTxValue struct {
	tx Transaction
}

type ctxReadonlyValue struct {
	db Readonly
}

func getTxFromContext(ctx context.Context) (ctxTxValue, bool) {
	tx, ok := ctx.Value(ctxTxKey).(ctxTxValue)
	return tx, ok
}

// GetTx get Transaction from context
func GetTx(ctx context.Context) Transaction {
	tx, ok := getTxFromContext(ctx)
	if !ok {
		panic("Not found transaction object in context")
	}
	return tx.tx
}

// GetReadonly get Readonly from context
func GetReadonly(ctx context.Context) Readonly {
	db, ok := ctx.Value(ctxReadonlyKey).(ctxReadonlyValue)
	if ok {
		return db.db
	}

	tx, ok := getTxFromContext(ctx)
	if ok {
		return tx.tx
	}

	panic("Not found readonly object in context")
}
