package bnats

const (
	prefix = "[nats] "

	errStreamCtx   = prefix + "ошибка создания конекста стрима: %w"
	errAddStream   = prefix + "ошибка создания стрима: %w"
	errInfoStream  = prefix + "ошибка получения статуса стрима: %w"
	errSubscribe   = prefix + "ошибка создания подписчика: %w"
	errUnsubscribe = prefix + "ошибка закрытия подписчика nats: %s"

	errPush  = prefix + "ошибка отправки сообщения в топик %s: %w"
	errFetch = prefix + "ошибка получения сообщений из топика %s: %w"

	errInsertLog = prefix + "ошибка записи лога: %w"
)
