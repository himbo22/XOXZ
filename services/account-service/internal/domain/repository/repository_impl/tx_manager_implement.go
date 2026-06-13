package repository_impl

import (
	"context"

	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) repository.TxRunner {
	return &TxManager{db: db}
}

// RunInTx implements repository.TxRunner.
func (t *TxManager) RunInTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := t.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	txCtx := context.WithValue(ctx, _const.TxKey, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(_const.TxKey).(*gorm.DB)
	if !ok || tx == nil {
		return db.WithContext(ctx)
	}
	return tx
}
