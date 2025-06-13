package rredis

const (
	errRead      = "ошибка чтения redis: %w"
	errCreateKey = "ошибка записи ключа %s redis: %w"
	errDeleteKey = "ошибка удаления ключа %s redis: %w"
)
