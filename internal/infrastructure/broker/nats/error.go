package bnats

const (
	errStreamCtx  = "ошибка создания конекста стрима: %w"
	errAddStream  = "ошибка создания стрима: %w"
	errInfoStream = "ошибка получения статуса стрима: %w"
	errSubscribe  = "ошибка создания подписчика: %w"

	errPush  = "ошибка отправки сообщения в топик %s: %w"
	errFetch = "ошибка получения сообщений из топика %s: %w"

	errInsertLog = "ошибка записи лога: %w"
)
