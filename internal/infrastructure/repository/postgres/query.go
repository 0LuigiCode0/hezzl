package rpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func query[T any](db _IConn, ctx context.Context, query string, arg ...any) (*T, error) {
	rows, err := db.Query(ctx, query, arg...)
	if err != nil {
		return nil, fmt.Errorf(errSelect, err)
	}
	out, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(errScanRow, err)
	}

	return out, nil
}

func queryRows[T any](db _IConn, ctx context.Context, query string, arg ...any) ([]*T, error) {
	rows, err := db.Query(ctx, query, arg...)
	if err != nil {
		return nil, fmt.Errorf(errSelect, err)
	}
	out, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(errScanRow, err)
	}

	return out, nil
}
