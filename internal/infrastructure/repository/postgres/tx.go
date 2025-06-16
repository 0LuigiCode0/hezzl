package rpostgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type _tx struct {
	_repo
	tx pgx.Tx
}

//go:generate mockery --name ITx --outpkg mpostgres --recursive --with-expecter
type ITx interface {
	_IRepo

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	EndTx(ctx context.Context, err error)
}

func (pg *_postgres) Begin(ctx context.Context) (ITx, error) {
	tx, err := pg.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(errTxCreate, err)
	}

	return &_tx{_repo: _repo{tx}, tx: tx}, nil
}

func (tx *_tx) Commit(ctx context.Context) error {
	err := tx.tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf(errTxCommit, err)
	}

	return nil
}

func (tx *_tx) Rollback(ctx context.Context) error {
	err := tx.tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf(errTxRollback, err)
	}

	return nil
}

func (tx *_tx) EndTx(ctx context.Context, err error) {
	if err != nil {
		err = tx.tx.Rollback(ctx)
		if err != nil {
			log.Print(fmt.Errorf(errTxRollback, err))
		}
	} else {
		err = tx.tx.Commit(ctx)
		if err != nil {
			log.Print(fmt.Errorf(errTxCommit, err))
		}
	}
}
