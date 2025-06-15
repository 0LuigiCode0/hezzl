package rpostgres

const (
	errTxCreate   = "ошибка создания транзакции: %w"
	errTxCommit   = "ошибка коммита транзакции: %w"
	errTxRollback = "ошибка отката транзакции: %w"

	errSelect  = "ошибка select: %w"
	errScanRow = "ошибка скана строки: %w"

	errInsertGood   = "ошибка записи товара: %w"
	errUpdateGood   = "ошибка обновления товара: %w"
	errRemoveGood   = "ошибка удаления товара: %w"
	errGetGoods     = "ошибка поиска товаров: %w"
	errGetGoodsMeta = "ошибка поиска статистики товаров: %w"
)
