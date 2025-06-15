package rpostgres

const (
	prefix = "[postgres] "

	errTxCreate   = prefix + "ошибка создания транзакции: %w"
	errTxCommit   = prefix + "ошибка коммита транзакции: %w"
	errTxRollback = prefix + "ошибка отката транзакции: %w"

	errSelect  = prefix + "ошибка select: %w"
	errScanRow = prefix + "ошибка скана строки: %w"

	errInsertGood   = prefix + "ошибка записи товара: %w"
	errUpdateGood   = prefix + "ошибка обновления товара: %w"
	errRemoveGood   = prefix + "ошибка удаления товара: %w"
	errGetGoods     = prefix + "ошибка поиска товаров: %w"
	errGetGoodsMeta = prefix + "ошибка поиска статистики товаров: %w"
)
