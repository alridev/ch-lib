# ch-lib

Библиотека для работы с CH API на Go.

## Основные возможности

- Отправка одиночных событий в CH
- Пакетная отправка событий
- Отправка данных в конкретную таблицу

## Примеры использования

### Установка модуля
```bash
go get github.com/alridev/ch-lib
```

### Создание клиента

```go

import (
    "github.com/alridev/ch-lib"
)

client := chlib.NewChClient(
    "http://localhost:8080",
    "setter-token",
    "getter-token",
    nil,
)
```

### Отправка одного события

```go
err := client.SetterOne(ctx, chlib.SetterEvent{
    TableName: "events",
    Data:      map[string]interface{}{"field": "value"},
})
```

### Пакетная отправка

```go
errors, err := client.SetterMany(ctx, []chlib.SetterEvent{
    {
        TableName: "events",
        Data:      map[string]interface{}{"field": "value"},
    },
})
```

### Отправка в конкретную таблицу

```go
err = client.SetterByTable(ctx, "events", map[string]interface{}{"field": "value"})
```

## Структуры данных

### SetterEvent

```go
type SetterEvent struct {
    TableName string      `json:"table_name" binding:"required"` // Имя таблицы в CH
    Data      interface{} `json:"data" binding:"required"`       // Данные для отправки
}
```

### SetterRequestBody

```go
type SetterRequestBody struct {
    Events []SetterEvent `json:"events" binding:"required"` // Список событий для отправки
}
```

### Ответы API

```go
// Ответ на пакетную отправку
type SetterManyResponseBody struct {
    Errors map[int32]string `json:"errors"` // Ошибки по индексам событий
}

// Ответ на отправку в конкретную таблицу
type SetterByTableResponseBody struct {
    Error string `json:"error"` // Текст ошибки
}

// Ответ на отправку одного события
type SetterOneResponseBody struct {
    Error string `json:"error"` // Текст ошибки
}
```
