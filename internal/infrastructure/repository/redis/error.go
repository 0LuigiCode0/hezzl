package rredis

const (
	prefix = "[redis] "

	errTx     = prefix + "ошибка выполнения транзакции: %w"
	errExpire = prefix + "ошибка установки срока ключа %s: %w"

	errReadKey   = prefix + "ошибка чтения ключа %s: %w"
	errCreateKey = prefix + "ошибка записи ключа %s: %w"
	errDeleteKey = prefix + "ошибка удаления ключа %s: %w"

	errAddSet = prefix + "ошибка добавления в список %s: %w"
	errRemSet = prefix + "ошибка удаления из списка %s: %w"
)
