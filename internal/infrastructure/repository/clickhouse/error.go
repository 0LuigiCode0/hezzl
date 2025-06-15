package rclickhouse

const (
	prefix = "[clickhouse] "

	errBatch       = prefix + "ошибка создания батча: %w"
	errBatchAppend = prefix + "ошибка добавления строки батча: %w"
	errBatchSend   = prefix + "ошибка отправки батча: %w"

	errInsertGoodsLog = prefix + "ошибка записи лога товаров: %w"
)
