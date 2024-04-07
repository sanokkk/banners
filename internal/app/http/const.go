package http

import "errors"

var (
	ErrNotAuthorized = errors.New("Пользователь не авторизован")
	ErrForbidden     = errors.New("Пользователь не имеет доступа")
	ErrTimeout       = errors.New("Время ответа сервера истекло")
)
