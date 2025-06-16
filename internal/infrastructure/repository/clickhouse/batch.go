package rclickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type _batch struct {
	b driver.Batch
}

type IBatch interface {
	AppendStruct(v IPrepareBatch) error
	Append(v ...any) error
	Send() error
}

type IPrepareBatch interface {
	PrepareBatch() []any
}

func (ch *_clickhouse) createBatch(ctx context.Context, query string) (IBatch, error) {
	b, err := ch.conn.PrepareBatch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(errBatch, err)
	}
	return &_batch{b: b}, nil
}

func (b *_batch) Send() error {
	err := b.b.Send()
	if err != nil {
		return fmt.Errorf(errBatchSend, err)
	}

	return nil
}

func (b *_batch) Append(v ...any) error {
	err := b.b.Append(v)
	if err != nil {
		return fmt.Errorf(errBatchAppend, err)
	}

	return nil
}

func (b *_batch) AppendStruct(v IPrepareBatch) error {
	err := b.b.Append(v.PrepareBatch()...)
	if err != nil {
		return fmt.Errorf(errBatchAppend, err)
	}

	return nil
}
