package service

import (
	"context"
	"github.com/nqxcode/platform_common/client/db"
)

type txManagerFake struct {
}

func NewTxManagerFake() db.TxManager {
	return &txManagerFake{}
}

func (tx *txManagerFake) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}
