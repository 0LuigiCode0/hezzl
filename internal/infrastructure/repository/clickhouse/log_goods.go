package rclickhouse

import (
	"context"
	"fmt"
)

const (
	qInsertGoodsLog = "insert into log_goods_store ( `id`,`project_id`,`name`,`description`,`priority`,`removed`,`eventTime`)"
)

func (ch *_clickhouse) InsertGoodsLogBatch(ctx context.Context) (IBatch, error) {
	res, err := ch.createBatch(ctx, qInsertGoodsLog)
	if err != nil {
		return nil, fmt.Errorf(errInsertGoodsLog, err)
	}

	return res, nil
}
