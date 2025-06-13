package rpostgres

const (
	errTxCreate   = "ошибка создания транзакции: %w"
	errTxCommit   = "ошибка коммита транзакции: %w"
	errTxRollback = "ошибка отката транзакции: %w"

	errSelect  = "ошибка select: %w"
	errScanRow = "ошибка скана строки: %w"

	errExec = "ошибка выполнения команды: %w"
)
