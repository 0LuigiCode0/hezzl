package ihttp

const (
	prefix = "[router] "

	errNotFound   = prefix + "errors.common.notFound"
	errParseBody  = prefix + "ошибка парса тела: %s"
	errParseParam = prefix + "ошибка парса параметра: %s"

	errReadBody = prefix + "ошибка чтения тела: %w"
)
