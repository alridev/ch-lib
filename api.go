package chlib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SetterResponseBody представляет базовую структуру ответа с ошибками
type SetterResponseBody struct {
	Errors map[int32]string `json:"errors"` // Ошибки по индексам
}

// ChEndpoints содержит пути к эндпоинтам API
type ChEndpoints struct {
	SetterMany    string // Путь для пакетной отправки
	SetterByTable string // Путь для отправки в таблицу (%s заменяется на имя таблицы)
	SetterOne     string // Путь для отправки одного события
}

// DefaultEndpoints содержит стандартные пути к эндпоинтам
var DefaultEndpoints = ChEndpoints{
	SetterMany:    "/setter/many",
	SetterByTable: "/setter/%s", // %s заменяется на имя таблицы
	SetterOne:     "/setter",
}

// ChClient представляет клиент для работы с CH API
type ChClient struct {
	BaseURL string // Базовый URL API

	SetterToken string // Токен для API записи
	GetterToken string // Токен для API чтения

	Endpoints ChEndpoints // Настройки путей API (опционально)
}

// NewChClient создает новый экземпляр клиента
// baseURL - базовый URL API
// setterToken - токен для API записи
// getterToken - токен для API чтения
// endpoints - настройки путей API (если nil, используются DefaultEndpoints)
func NewChClient(baseURL, setterToken, getterToken string, endpoints *ChEndpoints) *ChClient {
	if endpoints == nil {
		endpoints = &DefaultEndpoints
	}

	return &ChClient{
		BaseURL:     baseURL,
		SetterToken: setterToken,
		GetterToken: getterToken,
		Endpoints:   *endpoints,
	}
}

// SetterMany отправляет множество событий в CH
func (c *ChClient) SetterMany(ctx context.Context, events []SetterEvent) (map[int32]string, error) {
	url := c.BaseURL + c.Endpoints.SetterMany

	requestBody := SetterRequestBody{
		Events: events,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.SetterToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response SetterManyResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return response.Errors, nil
	}

	var response SetterManyResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Errors, nil
}

// SetterOne отправляет одно событие в CH
func (c *ChClient) SetterOne(ctx context.Context, event SetterEvent) error {
	url := c.BaseURL + c.Endpoints.SetterOne

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.SetterToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response SetterOneResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// SetterByTable отправляет данные в указанную таблицу
func (c *ChClient) SetterByTable(ctx context.Context, tableName string, data interface{}) error {
	url := c.BaseURL + fmt.Sprintf(c.Endpoints.SetterByTable, tableName)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.SetterToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response SetterByTableResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}
