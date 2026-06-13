package repository

import "context"

type TxRunner interface {
	RunInTx(ctx context.Context, fn func(txCtx context.Context) error) error
}
