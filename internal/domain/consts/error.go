package consts

const (
	ErrOpenConnect  = "ошибка подключения: %w"
	ErrCloseConnect = "ошибка закрытия соединения: %s"

	ErrJsonMarshal   = "ошибка чтения json: %w"
	ErrJsonUnmarshal = "ошибка записи json: %w"

	ErrPing = "ошибка пинга: %w"
)

const (
	ErrFieldValid = "поле %s невалидно"
)

const (
	NotifyClose = "закрыт"
)
