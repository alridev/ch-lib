package chlib

// SetterEvent представляет событие для отправки в CH
type SetterEvent struct {
	TableName string      `json:"table_name" binding:"required"` // Имя таблицы в CH
	Data      interface{} `json:"data" binding:"required"`       // Данные для отправки
}

// SetterRequestBody представляет тело запроса для пакетной отправки
type SetterRequestBody struct {
	Events []SetterEvent `json:"events" binding:"required"` // Список событий для отправки
}

// SetterManyResponseBody представляет ответ на пакетную отправку
type SetterManyResponseBody struct {
	Errors map[int32]string `json:"errors"` // Ошибки по индексам событий
}

// SetterByTableResponseBody представляет ответ на отправку в конкретную таблицу
type SetterByTableResponseBody struct {
	Error string `json:"error"` // Текст ошибки
}

// SetterOneResponseBody представляет ответ на отправку одного события
type SetterOneResponseBody struct {
	Error string `json:"error"` // Текст ошибки
}

