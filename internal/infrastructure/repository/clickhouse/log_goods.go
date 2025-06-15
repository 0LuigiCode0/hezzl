package rclickhouse

import (
	"context"
)

const (
	qInsertGoodsLog = "insert into log_goods_store ( `id`,`project_id`,`name`,`description`,`priority`,`removed`,`eventTime`)"
)

func (ch *_clickhouse) InsertGoodsLogBatch(ctx context.Context) (IBatch, error) {
	return ch.createBatch(ctx, qInsertGoodsLog)
}
